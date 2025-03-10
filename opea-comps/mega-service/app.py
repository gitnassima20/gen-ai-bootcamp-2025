from fastapi import HTTPException
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo
)
from comps.cores.mega.constants import ServiceType, ServiceRoleType
from comps import MicroService, ServiceOrchestrator
import os
import time

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
    async def handle_request(self, request: ChatCompletionRequest) -> ChatCompletionResponse:
        try:
            # Format the request for Ollama
            
            ollama_request = {
                "model": request.model or "llama3.2:1b",  
                "messages": request.messages,
                "stream": False  
            }
            
            print("\n\n\n Ollama Req", ollama_request)
            
            # Schedule the request through the orchestrator
            result = await self.megaservice.schedule(ollama_request)
            print("\n\n\nResult", result)
            
            # Extract the actual content from the response
            if isinstance(result, tuple) and len(result) > 0:
                llm_response = result[0].get('llm/MicroService')
                print("\n\n\nLLM Response", llm_response)
                
                # Handle StreamingResponse
                if hasattr(llm_response, 'json'):
                    # Read and process the response
                    response_body = await llm_response.json()
                    print("\n\n\nDecoded Content:", response_body)
                else:
                    content = "No response content available"

                # Construct and return the response
                return ChatCompletionResponse(
                    id="chatcmpl-" + str(time.time()),
                    object="chat.completion",
                    created=int(time.time()),
                    model=request.model or "llama3.2:1b",
                    choices=[
                        ChatCompletionResponseChoice(
                            index=0,
                            message=ChatMessage(
                                role="assistant",
                                content=content
                            ),
                            finish_reason="stop"
                        )
                    ],
                    usage=UsageInfo(
                        prompt_tokens=len(request.messages[0]['content'].split()) if request.messages else 0,
                        completion_tokens=len(content.split()),
                        total_tokens=len(request.messages[0]['content'].split()) + len(content.split()) if request.messages else 0
                    )
                )
            else:
                raise HTTPException(status_code=500, detail="No response from LLM service")

        except Exception as e:
            print(f"Error in handle_request: {e}")
            import traceback
            traceback.print_exc()
            raise HTTPException(status_code=500, detail=f"Internal Server Error: {str(e)}")

example = ExampleService()
example.add_remote_service()
example.start()