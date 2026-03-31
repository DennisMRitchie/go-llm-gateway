package gateway

import (
	"context"
	"fmt"
	"log"
)

// FallbackChain tries each provider in order until one succeeds.
type FallbackChain struct {
	providers []Provider
}

// NewFallbackChain builds a chain from an ordered slice of providers.
func NewFallbackChain(providers ...Provider) *FallbackChain {
	return &FallbackChain{providers: providers}
}

// Execute runs the chain, returning the first successful response.
func (fc *FallbackChain) Execute(ctx context.Context, req ChatRequest) (ChatResponse, string, error) {
	var lastErr error

	for _, p := range fc.providers {
		if !p.Healthy() {
			log.Printf("[fallback] provider=%s is unhealthy, skipping", p.Name())
			continue
		}

		resp, err := p.Generate(ctx, req)
		if err == nil {
			return resp, p.Name(), nil
		}

		log.Printf("[fallback] provider=%s failed: %v", p.Name(), err)
		lastErr = err
	}

	return ChatResponse{}, "", fmt.Errorf("all providers failed; last error: %w", lastErr)
}
