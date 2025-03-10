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
    "max_tokens": 100
  }'\
  -o output/response.json
'
```
