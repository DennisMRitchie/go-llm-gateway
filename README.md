# 🚀 go-llm-gateway

**High-Performance OpenAI-Compatible LLM Gateway written in pure Go**

A unified gateway for multiple LLM providers (Ollama, OpenAI, Groq, and more) with intelligent routing, automatic fallback, rate limiting, and observability — ready for production and Kubernetes.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![OpenAI Compatible](https://img.shields.io/badge/OpenAI_Compatible-Yes-10A37F?style=flat)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-yellow?style=flat)

---

## ✨ Key Features

| Feature | Details |
|---|---|
| ⚡ Drop-in OpenAI replacement | `/v1/chat/completions` endpoint |
| 🧠 Intelligent routing | Routes by model keyword (qwen, deepseek, gpt, llama) |
| 🔄 Automatic fallback | Cascades through providers on error |
| 🛡 Per-IP rate limiting | Token-bucket limiter via `golang.org/x/time/rate` |
| 📊 Observability | Request tracing + `/metrics` endpoint |
| 🔗 Provider extensible | Add any provider by implementing `Provider` interface |
| 🐳 Docker / Kubernetes ready | Multi-stage Dockerfile, compose included |

---

## 🏗 Architecture

```
Client (any OpenAI SDK)
        │
        ▼
  ┌─────────────┐
  │  Gin HTTP   │  rate limit · tracing
  └──────┬──────┘
         │
  ┌──────▼──────┐
  │   Router    │  keyword-based model routing
  └──────┬──────┘
         │
  ┌──────▼──────┐
  │  Provider   │  Ollama · OpenAI · LocalFallback
  └─────────────┘
```

---

## 🚀 Quick Start

```bash
git clone https://github.com/DennisMRitchie/go-llm-gateway.git
cd go-llm-gateway
make up       # build & start Docker services
make test     # send a test request
make health   # check service health
```

---

## 📡 API

### `POST /v1/chat/completions`

OpenAI-compatible chat endpoint.

```bash
curl -X POST http://localhost:8082/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "qwen3",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

**Response:**
```json
{
  "response": {
    "id": "chat-20260331120000",
    "object": "chat.completion",
    "model": "qwen3",
    "choices": [{
      "index": 0,
      "message": {"role": "assistant", "content": "..."},
      "finish_reason": "stop"
    }]
  },
  "used_provider": "ollama",
  "gateway_version": "v0.3.0"
}
```

### `GET /health`
Returns service status and total request count.

### `GET /metrics`
Returns basic metrics (requests total).

---

## ⚙️ Configuration

All settings via environment variables:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8082` | Server port |
| `OLLAMA_URL` | `http://ollama:11434` | Ollama base URL |
| `OPENAI_API_KEY` | `` | OpenAI API key (optional) |
| `DEFAULT_MODEL` | `qwen3` | Model used when none specified |

---

## 🧩 Routing Rules

| Model keyword | Provider |
|---|---|
| `qwen*` | ollama |
| `deepseek*` | ollama |
| `llama*` | ollama |
| `gpt*` | openai |
| *(default)* | ollama → local-fallback |

---

## 🛠 Tech Stack

- **Go 1.23** + Gin
- Multi-provider abstraction (`Provider` interface)
- Per-IP rate limiting (`golang.org/x/time/rate`)
- OpenTelemetry-ready tracing middleware
- Docker Compose + multi-stage Dockerfile

---

## 🔗 Related Projects

- [`go-rag-llm-orchestrator`](https://github.com/DennisMRitchie/go-rag-llm-orchestrator) — RAG pipeline
- [`go-llm-smart-cache`](https://github.com/DennisMRitchie/go-llm-smart-cache) — Semantic caching layer

Together these form a production-ready **Go LLM stack**.

---

## 📄 License

MIT

---

Built with ❤️ by **Konstantin Lychkov**  
Senior Go Developer | Go + LLM/NLP Specialist  
Open to Remote Worldwide

⭐ If you find this useful, give it a star!
