
``` sh
curl -X POST http://localhost:8000/v1/example-service \
  -H "Content-Type: application/json" \
  -d '{
    "messages": "Hello, this is a test message, please respond"
  }'\
  -o response.json
'
```