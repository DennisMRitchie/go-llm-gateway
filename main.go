package main

import (
	"context"
	"log"
	"net/http"

	"github.com/DennisMRitchie/go-llm-gateway/internal/config"
	"github.com/DennisMRitchie/go-llm-gateway/internal/gateway"
	"github.com/DennisMRitchie/go-llm-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// ── Providers ────────────────────────────────────────────────────────────
	router := gateway.NewRouter()
	router.Register("ollama", gateway.NewOllamaProvider(cfg.OllamaURL))
	router.Register("local-fallback", gateway.NewLocalFallbackProvider())

	if cfg.OpenAIAPIKey != "" {
		router.Register("openai", gateway.NewOpenAIProvider(cfg.OpenAIAPIKey))
		log.Println("[main] OpenAI provider registered")
	}

	// ── HTTP server ───────────────────────────────────────────────────────────
	r := gin.Default()
	r.Use(middleware.Tracing())
	r.Use(middleware.RateLimit())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "v0.3.0",
			"requests_total": router.RequestCount(),
		})
	})

	// OpenAI-compatible chat completions endpoint
	r.POST("/v1/chat/completions", func(c *gin.Context) {
		var req gateway.ChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Model == "" {
			req.Model = cfg.DefaultModel
		}

		resp, usedProvider, err := router.Route(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"response":        resp,
			"used_provider":   usedProvider,
			"gateway_version": "v0.3.0",
		})
	})

	// Metrics endpoint (simple)
	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"requests_total": router.RequestCount(),
		})
	})

	log.Printf("🚀 Go LLM Gateway running on :%s | OpenAI-compatible endpoint ready", cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
