# Running TTS Service

## Containers

- `speecht5-service` performs raw speech tasks but lacks a public interface.

- `tts` provides a general TTS service layer that handles request management.

- `tts-speecht5` specializes the TTS layer to route requests specifically to SpeechT5.

## Making a Request to the server

``` sh
curl http://localhost:9088/v1/audio/speech -XPOST \
-d '{"input":"Who are you?"}' \
-H 'Content-Type: application/json' \
--output speech.wav

```

## Request for GPTSoVITS

``` sh
curl -X POST "http://:9880 \
-H "Content-Type: application/json" \
-d ' {
     "prompt_text":"Hellow world",
     "prompt_language":"en",
     "text": "This is a new sentence I want to convert it to speech",
     "text_language":"en",
}'

```

## Problem with GPTSoVITS

- So large image about 23GB
- Requires a lot of memory allocation
