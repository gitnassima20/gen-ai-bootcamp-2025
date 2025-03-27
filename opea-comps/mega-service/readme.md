# Guide

## How to run the LLM service

Setting a port to `9000` since it's the ideal one when it comes to
OPEA microservices default ports
It will default to `8008` as in the docker-compose file if not set.

``` sh
LLM_ENDPOINT_PORT=9000 docker-compose up
```

## How to run the MegaService

``` sh
python app.py
```

## How to access the Jaeger UI

After running the docker-compose the jaeger service will be available at `http://localhost:16686/search`

## How to make a request to the FastAPI

``` sh
curl -X POST http://localhost:8000/v1/example-service \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [ 
      {
        "role": "user",
        "content": "Hello, this is a test message, please respond"
      }
    ],
    "model": "llama3.2:1b",
    "temperature": 0.7,
    "max_tokens": 100,
    "stream": false
  }'\
  -o output/res.json
'
```

## Resolving error

### error:-`{'llm/MicroService': {'error': {'message': "[] is too short - 'messages'", 'type': 'invalid_request_error', 'param': None, 'code': None}}}`

We adjust the `initial_inputs` to include the original messages

```python
initial_inputs={
                "text": prompt,
                "messages": chat_request.messages
            },
```

### error:-`result_dict {'llm/MicroService': {'error': {'message': 'model is required', 'type': 'api_error', 'param': None, 'code': None}}}`

the LLMParams were not passing the model

``` python
parameters = LLMParams(
            model=chat_request.model,
            ...)
```

## Test the overall combination of OPEA Microservices

``` sh
curl -X POST http://localhost:8000/v1/example-service   -H "Content-Type: application/json"   -d '{
    "messages": [{
      "role": "user",
      "content": "Hello, this is a test message"
    }],
    "model": "llama3.2:1b",
    "temperature": 0.7,
    "max_tokens": 100
  }'   --create-dirs   -o output/res.wav
```