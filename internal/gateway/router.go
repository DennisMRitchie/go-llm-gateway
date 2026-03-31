package gateway

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
)

// modelRouteRule maps a keyword in the model name to a preferred provider name.
type modelRouteRule struct {
	keyword  string
	provider string
}

// Router performs intelligent routing across registered providers.
type Router struct {
	providers       map[string]Provider
	defaultProvider string
	fallbackOrder   []string
	rules           []modelRouteRule
	requests        uint64 // atomic counter
}

// NewRouter returns a Router with sensible defaults.
func NewRouter() *Router {
	return &Router{
		providers:       make(map[string]Provider),
		defaultProvider: "ollama",
		fallbackOrder:   []string{"ollama", "local-fallback"},
		rules: []modelRouteRule{
			{keyword: "qwen", provider: "ollama"},
			{keyword: "deepseek", provider: "ollama"},
			{keyword: "llama", provider: "ollama"},
			{keyword: "gpt", provider: "openai"},
		},
	}
}

// Register adds a provider to the router under the given name.
func (r *Router) Register(name string, p Provider) {
	r.providers[name] = p
	log.Printf("[router] registered provider=%s", name)
}

// RequestCount returns total requests handled so far.
func (r *Router) RequestCount() uint64 {
	return atomic.LoadUint64(&r.requests)
}

// Route picks the best provider for the request and executes it.
// Falls back through r.fallbackOrder on error.
func (r *Router) Route(ctx context.Context, req ChatRequest) (ChatResponse, string, error) {
	atomic.AddUint64(&r.requests, 1)

	providerName := r.selectProvider(req.Model)
	log.Printf("[router] request=%d model=%s → provider=%s", r.RequestCount(), req.Model, providerName)

	p, ok := r.providers[providerName]
	if !ok {
		log.Printf("[router] provider=%s not registered, falling back", providerName)
		return r.runFallback(ctx, req)
	}

	resp, err := p.Generate(ctx, req)
	if err != nil {
		log.Printf("[router] provider=%s error: %v — starting fallback", providerName, err)
		return r.runFallback(ctx, req)
	}

	return resp, providerName, nil
}

// selectProvider returns the provider name based on routing rules.
func (r *Router) selectProvider(model string) string {
	lower := strings.ToLower(model)
	for _, rule := range r.rules {
		if strings.Contains(lower, rule.keyword) {
			return rule.provider
		}
	}
	return r.defaultProvider
}

// runFallback executes the fallback chain.
func (r *Router) runFallback(ctx context.Context, req ChatRequest) (ChatResponse, string, error) {
	var chain []Provider
	for _, name := range r.fallbackOrder {
		if p, ok := r.providers[name]; ok {
			chain = append(chain, p)
		}
	}
	return NewFallbackChain(chain...).Execute(ctx, req)
}
