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

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 6000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "localhost")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 9000)


class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        self.endpoint = "/v1/example-service"
        self.megaservice = ServiceOrchestrator()

    def add_remote_service(self):
        #embedding = MicroService(
        #    name="embedding",
        #    host=EMBEDDING_SERVICE_HOST_IP,
        #    port=EMBEDDING_SERVICE_PORT,
        #    endpoint="/v1/embeddings",
        #    use_remote_service=True,
        #    service_type=ServiceType.EMBEDDING,
        #)
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        #self.megaservice.add(embedding).add(llm)
        #self.megaservice.flow_to(embedding, llm)
        self.megaservice.add(llm)
    
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
        # for node, response in result_dict.items():
        #     if isinstance(response, StreamingResponse):
        #         return response
        # last_node = runtime_graph.all_leaves()[-1]
        # print("\n\nlast_node",last_node)
        # response = result_dict[last_node]["text"]
        # print("\n\nresponse",response)
        
        # try:
        #     # Comprehensive request validation and logging
        #     print("\n\n--- REQUEST VALIDATION ---")
        #     print(f"Request Type: {type(request)}")
        #     print(f"Request Attributes: {request}")
            
        #     # Validate request structure
        #     if not hasattr(request, 'messages') or not request.messages:
        #         print("ERROR: No messages found in the request")
        #         raise HTTPException(status_code=400, detail="Messages are required")
            
        #     # Validate messages structure
        #     if not isinstance(request.messages, list):
        #         print(f"ERROR: Messages must be a list, got {type(request.messages)}")
        #         raise HTTPException(status_code=400, detail="Messages must be a list")
            
        #     # Validate each message
        #     for msg in request.messages:
        #         if not isinstance(msg, dict):
        #             print(f"ERROR: Invalid message format: {msg}")
        #             raise HTTPException(status_code=400, detail="Each message must be a dictionary")
                
        #         if 'role' not in msg or 'content' not in msg:
        #             print(f"ERROR: Message missing required keys: {msg}")
        #             raise HTTPException(status_code=400, detail="Each message must have 'role' and 'content' keys")
            
        #     # Prepare Ollama request
        #     ollama_request = {
        #         "model": request.model or "llama3.2:1b",
        #         "messages": request.messages,
        #         "temperature": request.temperature if hasattr(request, 'temperature') else 0.7,
        #         "stream": False
        #     }
            
        #     # Optional: Add additional parameters if present
        #     if hasattr(request, 'max_tokens') and request.max_tokens is not None:
        #         ollama_request["max_tokens"] = request.max_tokens
            
        #     print("\n\n--- OLLAMA REQUEST ---")
        #     print(ollama_request)
            
        #     # Schedule the request through the orchestrator
        #     result = await self.megaservice.schedule(ollama_request)
        #     print("\n\n--- ORCHESTRATOR RESULT ---")
        #     print(result)
            
        #     # Extract the actual content from the response
        #     if isinstance(result, tuple) and len(result) > 0:
        #         llm_response = result[0].get('llm/MicroService')
        #         print("\n\n--- LLM RESPONSE ---")
        #         print(llm_response)
                
        #         # Handle StreamingResponse
        #         if hasattr(llm_response, 'body_iterator'):
        #             # Read and process the response
        #             response_body = b""
        #             async for chunk in llm_response.body_iterator:
        #                 response_body += chunk
        #             content = response_body.decode('utf-8')
        #             print("\n\n--- DECODED CONTENT ---")
        #             print(content)
        #         elif hasattr(llm_response, 'body'):
        #             # Alternative method to extract content
        #             response_body = await llm_response.body()
        #             content = response_body.decode('utf-8')
        #             print("\n\n--- DECODED CONTENT (body method) ---")
        #             print(content)
        #         else:
        #             content = "No response content available"
        #             print("\n\n--- NO CONTENT FOUND ---")

        #         # Construct and return the response
        #         return ChatCompletionResponse(
        #             id="chatcmpl-" + str(time.time()),
        #             object="chat.completion",
        #             created=int(time.time()),
        #             model=request.model or "llama3.2:1b",
        #             choices=[
        #                 ChatCompletionResponseChoice(
        #                     index=0,
        #                     message=ChatMessage(
        #                         role="assistant",
        #                         content=content
        #                     ),
        #                     finish_reason="stop"
        #                 )
        #             ],
        #             usage=UsageInfo(
        #                 prompt_tokens=len(request.messages[0]['content'].split()) if request.messages else 0,
        #                 completion_tokens=len(content.split()),
        #                 total_tokens=len(request.messages[0]['content'].split()) + len(content.split()) if request.messages else 0
        #             )
        #         )
        #     else:
        #         raise HTTPException(status_code=500, detail="No response from LLM service")

        # except HTTPException:
        #     # Re-raise HTTP exceptions directly
        #     raise
        # except Exception as e:
        #     print(f"UNEXPECTED ERROR in handle_request: {e}")
        #     import traceback
        #     traceback.print_exc()
        #     raise HTTPException(status_code=500, detail=f"Internal Server Error: {str(e)}")

example = ExampleService()
example.add_remote_service()
example.start()