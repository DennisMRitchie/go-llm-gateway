# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/llm-gateway .

# ── Stage 2: minimal runtime image ───────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/llm-gateway .

EXPOSE 8082

CMD ["./llm-gateway"]
