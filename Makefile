.PHONY: up down build test health metrics tidy lint

## Start all services in background
up:
	docker compose up --build -d

## Stop all services
down:
	docker compose down

## Build the binary locally
build:
	go build -o bin/llm-gateway .

## Run tests
test-unit:
	go test ./...

## Send a test chat request
test:
	@curl -s -X POST http://localhost:8082/v1/chat/completions \
	  -H "Content-Type: application/json" \
	  -d '{ \
	    "model": "qwen3", \
	    "messages": [{"role": "user", "content": "What is the best open-source LLM in 2026?"}] \
	  }' | jq .

## Health check
health:
	@curl -s http://localhost:8082/health | jq .

## Simple metrics
metrics:
	@curl -s http://localhost:8082/metrics | jq .

## Tidy dependencies
tidy:
	go mod tidy

## Lint (requires golangci-lint)
lint:
	golangci-lint run ./...
