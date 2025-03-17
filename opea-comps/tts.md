# Running TTS Service

## Containers

- `speecht5-service` performs raw speech tasks but lacks a public interface.

- `tts` provides a general TTS service layer that handles request management.

- `tts-speecht5` specializes the TTS layer to route requests specifically to SpeechT5.
