from fastapi.responses import JSONResponse, StreamingResponse
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
)
from comps.cores.proto.docarray import LLMParams
from comps.cores.mega.utils import handle_message
from comps.cores.mega.constants import ServiceType, ServiceRoleType
from comps import MicroService, ServiceOrchestrator
from fastapi import Request
from fastapi.responses import StreamingResponse
import os
import logging
from io import BytesIO
import httpx

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

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
        self.generate_tts_response = self.generate_tts_response.__get__(self)

    def add_remote_service(self):
        tts = MicroService(
           name="tts",
           host=TTS_SERVICE_HOST_IP,
           port=TTS_SERVICE_PORT,
           endpoint="/v1/audio/speech",
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
        self.megaservice.add(llm).add(tts)
        self.megaservice.flow_to(llm, tts)
    
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
        try:
            data = await request.json()
            logger.info(f"Received request data: {data}")

            stream_opt = data.get("stream", False)
            chat_request = ChatCompletionRequest.model_validate(data)

            prompt = handle_message(chat_request.messages)

            parameters = LLMParams(
                model=chat_request.model,
                max_tokens=chat_request.max_tokens if chat_request.max_tokens else 1024,
                top_k=chat_request.top_k if chat_request.top_k else 10,
                top_p=chat_request.top_p if chat_request.top_p else 0.95,
                temperature=chat_request.temperature if chat_request.temperature else 0.01,
                stream=stream_opt,
            )

            logger.info(f"Processing prompt: {prompt}")

            result_dict, runtime_graph = await self.megaservice.schedule(
                initial_inputs={
                    "text": prompt,
                    "messages": chat_request.messages
                },
                llm_parameters=parameters
            )

            # Extract response text
            last_node = runtime_graph.all_leaves()[-1]
            node_response = result_dict[last_node]
            logger.info(f"Generated node response: {node_response}")

            # Add proper error handling and response extraction
            response_text = "Default response"  # Ensure a fallback value
            if isinstance(node_response, dict) and 'detail' in node_response:
                # Navigate through 'detail' to find 'choices'
                details = node_response['detail']
                if isinstance(details, list) and len(details) > 0:
                    input_data = details[0].get('input', {})  # Get 'input' safely
                    if isinstance(input_data, dict) and 'choices' in input_data:
                        response_text = input_data['choices'][0]['message']['content']

            logger.info(f"Generated response text: {response_text}")

            # TTS Service Request
            return await self.generate_tts_response(response_text, chat_request.model)

        except Exception as e:
            logger.error(f"Unexpected error: {str(e)}", exc_info=True)
            return JSONResponse(
                status_code=500,
                content={"error": f"Internal server error: {str(e)}"}
            )

    async def generate_tts_response(self, text: str, model: str):
        """
        Generate TTS response with async httpx client
        """
        tts_endpoint = f"http://{TTS_SERVICE_HOST_IP}:{TTS_SERVICE_PORT}/v1/audio/speech"

        try:
            logger.info(f"Attempting TTS with endpoint {tts_endpoint}")
            
            # Use async httpx client
            async with httpx.AsyncClient(timeout=30.0) as client:
                response = await client.post(
                    tts_endpoint,
                    json={"input": text},
                    headers={"Content-Type": "application/json"}
                )

            # Check response
            if response.status_code == 200:
                logger.info("Successfully generated TTS audio")
                
                # Create a bytes buffer to ensure streaming works
                audio_buffer = BytesIO(response.content)
                
                # Return streaming response
                return StreamingResponse(
                    audio_buffer,
                    media_type="audio/wav",
                    headers={
                        "Content-Disposition": 'attachment; filename="speech_output.wav"',
                        "Content-Length": str(len(response.content))
                    }
                )
            else:
                logger.error(f"TTS service returned status code {response.status_code}")
                logger.error(f"Response content: {response.text}")
                return JSONResponse(
                    status_code=response.status_code,
                    content={"error": f"TTS service error: {response.text}"}
                )

        except httpx.RequestError as e:
            logger.error(f"TTS request failed: {str(e)}")
            return JSONResponse(
                status_code=502,
                content={"error": f"TTS service connection error: {str(e)}"}
            )
        except Exception as e:
            logger.error(f"Unexpected TTS error: {str(e)}", exc_info=True)
            return JSONResponse(
                status_code=500,
                content={"error": "Could not generate speech from text"}
            )
example = ExampleService()
example.add_remote_service()
example.start()