# Running Ollama Third-Party Service

## Choosing a Model

You can get the model_id that ollama will launch from the [Ollama Library](https://ollama.com/library).

eg. <https://ollama.com/library/llama3.2:1b>

## Installing Ubuntu in top of WSL

Run in powershell as admin:

``` sh
curl.exe -L -o ubuntu-2204.appx https://aka.ms/wslubuntu2204
```

Then run:

``` sh
Add-AppxPackage .\ubuntu-2204.appx
```

## Getting the Host IP

### Linux

Get your IP address

``` sh
sudo apt install net-tools
ifconfig
```

Or you can try this way `$(hostname -I | awk '{print $1}')`

HOST_IP=$(hostname -I | awk '{print $1}') NO_PROXY=localhost LLM_ENDPOINT_PORT=8008 LLM_MODEL_ID="llama3.2:1b" docker compose up

### Ollama API

Once the Ollama server os running we can make a request to it:

<https://github.com/ollama/ollama/blob/main/docs/api.md>

## Download (Pull) the model

``` sh
curl --noproxy "*" http://localhost:8008/api/pull -d '{
  "model": "llama3.2:1b"
}'
```

## Generate a Request

``` sh
curl --noproxy "*" http://localhost:8008/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt":"Why is the sky blue?"
}'
```

# Technical Uncertainty

- Q1: Does brisge mode mean we can only access Ollama API with another model in the docker compose ?
- A1: No, the host machine can access it.

- Q2: Which port is being mapped 8008->11434 ?
- A2: 8008 is the port that the host machine will access, th eother is the guest port (the port of the service inside the container).

- Q3: If we pass the LLM_MODEL_ID to the Ollama server will it download the model when on start ?
- A3: It does not appear so. The ollama CLI might be running multiple API so you need to call the /pull api before trying to generate text.

- Q4: Will the mdoel be downloaded in the container ?
does that mean that the ml model will be deleted when the container is stopped ?
- A4: The model will be downloaded in the container. And vanish when the container is stopped. You need to mount a local drive to store the model, maybe more work needed to be done for this part.

- Q5: For LLM service which can do tex-generation it suggests it will only work TGI/vLLM and all you have to do is to have it running. Does TGI and vLLM have a standarized API or is there code to detect which one is running ? Do we have to use Xeon or Guadi processors ?
- A5: Yes, TGI and vLLM are similar to Ollama in that they are all LLM inference frameworks,all of them offer APIs with OpenAI compatibility, so in theory they should work the same.