.PHONY: up down build test health metrics tidy

up:
	docker compose up --build -d

down:
	docker compose down

build:
	go build -o bin/llm-gateway .

test:
	curl -s -X POST http://localhost:8082/v1/chat/completions \
	  -H "Content-Type: application/json" \
	  -d @testdata/request.json | jq .

health:
	curl -s http://localhost:8082/health | jq .

metrics:
	curl -s http://localhost:8082/metrics | jq .

tidy:
	go mod tidy
