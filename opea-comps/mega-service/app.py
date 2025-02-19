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
import sys
import time
import asyncio
import logging

# Aggressive OpenTelemetry and tracing disablement
import os
import sys

# Disable all tracing and telemetry environment variables
os.environ["TELEMETRY_ENDPOINT"] = ""
os.environ["OTEL_SDK_DISABLED"] = "true"
os.environ["OTEL_TRACES_EXPORTER"] = "none"
os.environ["OTEL_METRICS_EXPORTER"] = "none"
os.environ["OTEL_LOGS_EXPORTER"] = "none"
os.environ["OTEL_PYTHON_DISABLED"] = "1"
os.environ["OTEL_PYTHON_LOGGING_AUTO_INSTRUMENTATION_ENABLED"] = "false"

# Patch sys.modules to prevent OpenTelemetry initialization
sys.modules['opentelemetry'] = None
sys.modules['opentelemetry.trace'] = None
sys.modules['opentelemetry.sdk.trace'] = None
sys.modules['opentelemetry.exporter.otlp'] = None
sys.modules['opentelemetry.instrumentation'] = None

# Prevent any tracing or instrumentation
try:
    import opentelemetry
    opentelemetry.disable_tracing()
except (ImportError, AttributeError):
    pass

# Disable logging for specific modules
import logging
logging.getLogger('opentelemetry').setLevel(logging.CRITICAL)
logging.getLogger('urllib3').setLevel(logging.CRITICAL)
logging.getLogger('requests').setLevel(logging.CRITICAL)

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('megaservice.log')
    ]
)
logger = logging.getLogger(__name__)

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 6000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "localhost")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 9000)


class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):
        print('Initializing ExampleService')
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
            # Enhanced logging for request details
            logger.debug(f"Received raw request type: {type(request)}")
            logger.debug(f"Received request attributes: {vars(request) if hasattr(request, '__dict__') else request}")
            
            # Robust message content extraction
            if isinstance(request, str):
                message_content = request
            elif hasattr(request, 'messages'):
                if isinstance(request.messages, list) and request.messages:
                    message_content = request.messages[0].content if hasattr(request.messages[0], 'content') else str(request.messages[0])
                else:
                    message_content = "Hello"
            else:
                logger.error(f"Unrecognized request format: {request}")
                raise HTTPException(status_code=400, detail="Invalid request format")
            
            logger.info(f"Extracted message content: {message_content}")
            
            # Validate request before processing
            if not message_content:
                logger.warning("Empty message content received")
                raise HTTPException(status_code=400, detail="Message content cannot be empty")
            
            # Log detailed debug information
            logger.debug(f"Received request: {request}")
            
            # Log service connection attempt
            logger.info(f"Attempting to connect to LLM service at {LLM_SERVICE_HOST_IP}:{LLM_SERVICE_PORT}")
            
            # Fallback mechanism if MegaService fails
            try:
                # Format the request for Ollama
                ollama_request = {
                    "model": "llama3.2:1b",
                    "prompt": message_content
                }
                
                # Log routing attempt
                logger.debug("Routing request through MegaService")
                
                # Await the coroutine returned by schedule
                response = await self.megaservice.schedule(ollama_request)
                
                # Log raw response for debugging
                logger.debug(f"Raw MegaService response: {response}")
                
                # Extract response content
                if isinstance(response, tuple) and len(response) > 0:
                    # Try to extract response from the first element
                    first_response = response[0]
                    
                    # Log first response for debugging
                    logger.debug(f"First response: {first_response}")
                    
                    # Check if it's a dictionary with 'llm/MicroService'
                    if isinstance(first_response, dict) and 'llm/MicroService' in first_response:
                        llm_response = first_response['llm/MicroService']
                        
                        # Log LLM response for debugging
                        logger.debug(f"LLM Response type: {type(llm_response)}")
                        
                        # If it's a StreamingResponse, try to get content
                        if hasattr(llm_response, 'body_iterator'):
                            # Consume the body iterator
                            body_content = b''
                            async for chunk in llm_response.body_iterator:
                                body_content += chunk
                            response_text = body_content.decode('utf-8')
                        else:
                            response_text = str(llm_response)
                    else:
                        response_text = str(first_response)
                else:
                    response_text = str(response)
                
                # Create a dictionary with the response
                response = {"response": response_text}
            except Exception as inference_error:
                logger.error(f"MegaService routing failed: {inference_error}", exc_info=True)
                # Fallback to a simple response
                response = {"response": "Service is currently unavailable. Please try again later."}
            
            logger.info(f"Processed response: {response}")
            
            # Construct ChatCompletionResponse
            return ChatCompletionResponse(
                id="chatcmpl-123",
                object="chat.completion",
                created=int(time.time()),
                model="llama3.2:1b",
                choices=[
                    ChatCompletionResponseChoice(
                        index=0,
                        message=ChatMessage(
                            role="assistant", 
                            content=response.get('response', 'No response')
                        ),
                        finish_reason="stop"
                    )
                ],
                usage=UsageInfo(
                    prompt_tokens=len(message_content.split()) if message_content else 0,
                    completion_tokens=len(response.get('response', '').split()) if response else 0,
                    total_tokens=0
                )
            )
        except Exception as e:
            logger.error(f"Unexpected error in handle_request: {e}", exc_info=True)
            raise HTTPException(status_code=500, detail=str(e))

example = ExampleService()
example.add_remote_service()
example.start()