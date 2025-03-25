from fastapi import HTTPException
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo
)
from comps.cores.proto.docarray import LLMParams
from comps.cores.mega.utils import handle_message
from comps.cores.mega.constants import ServiceType, ServiceRoleType
from comps import MicroService, ServiceOrchestrator
from fastapi import Request
from fastapi.responses import StreamingResponse
import os

TTS_SERVICE_HOST_IP = os.getenv("TTS_SERVICE_HOST_IP", "localhost")
TTS_SERVICE_PORT = os.getenv("TTS_SERVICE_PORT", 9088)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "localhost")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 9000)


class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        self.endpoint = "/v1/example-service"
        self.megaservice = ServiceOrchestrator()

    def add_remote_service(self):
        tts = MicroService(
           name="tts",
           host=TTS_SERVICE_HOST_IP,
           port=TTS_SERVICE_PORT,
           endpoint="/v1/tts",
           use_remote_service=True,
           service_type=ServiceType.TTS,
        )
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        self.megaservice.add(tts).add(llm)
        self.megaservice.flow_to(tts, llm)
    
    def start(self):

        self.service = MicroService(
            self.__class__.__name__,
            service_role=ServiceRoleType.MEGASERVICE,
            host=self.host,
            port=self.port,
            endpoint=self.endpoint,
            input_datatype=ChatCompletionRequest,
            output_datatype=ChatCompletionResponse,
        )

        self.service.add_route(self.endpoint, self.handle_request, methods=["POST"])

        self.service.start()
    async def handle_request(self, request: Request):
        data = await request.json()
        print("data", data)
        stream_opt = data.get("stream", False)
        print(stream_opt)
        chat_request = ChatCompletionRequest.model_validate(data)
        print(chat_request)

        # Instead of just passing the prompt text, pass the entire structured data
        # that your microservice expects
        prompt = handle_message(chat_request.messages)

        parameters = LLMParams(
            model=chat_request.model,
            max_tokens=chat_request.max_tokens if chat_request.max_tokens else 1024,
            top_k=chat_request.top_k if chat_request.top_k else 10,
            top_p=chat_request.top_p if chat_request.top_p else 0.95,
            temperature=chat_request.temperature if chat_request.temperature else 0.01,
            frequency_penalty=chat_request.frequency_penalty if chat_request.frequency_penalty else 0.0,
            presence_penalty=chat_request.presence_penalty if chat_request.presence_penalty else 0.0,
            repetition_penalty=chat_request.repetition_penalty if chat_request.repetition_penalty else 1.03,
            stream=stream_opt,
            chat_template=chat_request.chat_template if chat_request.chat_template else None,
        )

        print("\n\nprompt", prompt, "\n\nparameters", parameters)

        # include both the original messages and the processed prompt
        result_dict, runtime_graph = await self.megaservice.schedule(
            initial_inputs={
                "text": prompt,
                "messages": chat_request.messages
            },
            llm_parameters=parameters
        )

        print("\n\nresult_dict", result_dict)
        for node, response in result_dict.items():
            if isinstance(response, StreamingResponse):
                return response

        last_node = runtime_graph.all_leaves()[-1]
        node_response = result_dict[last_node]

        # Add proper error handling and response extraction
        if "error" in node_response:
            # Handle error case - either return the error or a default message
            error_msg = node_response["error"]["message"]
            return JSONResponse(
                status_code=400,
                content={"error": error_msg}
            )
        elif "text" in node_response:
            response_text = node_response["text"]
        else:
            # Try to extract text from other potential keys or use a default response
            response_text = str(node_response)

        # Send response to TTS service for audio generation
        tts_url = f"http://{TTS_SERVICE_HOST_IP}:{TTS_SERVICE_PORT}/v1/tts"
        tts_response = requests.post(tts_url, json={"text": response_text})

        if tts_response.status_code == 200:
            # Assuming TTS service returns an audio file or URL
            audio_data = tts_response.json()
            return StreamingResponse(audio_data, media_type="audio/mpeg")
        else:
            return JSONResponse(
                status_code=500,
                content={"error": "Failed to generate audio"}
            )

        choices = []
        usage = UsageInfo()
        choices.append(
            ChatCompletionResponseChoice(
                index=0,
                message=ChatMessage(role="assistant", content=response_text),
                finish_reason="stop",
            )
        )
        return ChatCompletionResponse(model="chatqna", choices=choices, usage=usage)

example = ExampleService()
example.add_remote_service()
example.start()