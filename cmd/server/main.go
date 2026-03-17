package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lffm1994/rate-limiter/config"
	"github.com/lffm1994/rate-limiter/internal/handler"
	"github.com/lffm1994/rate-limiter/internal/limiter"
	ratemw "github.com/lffm1994/rate-limiter/internal/middleware"
)

func main() {
	cfg := config.Load()

	store := limiter.NewRedisStore(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPass,
		cfg.RedisDB,
	)

	l := limiter.NewLimiter(store, limiter.Config{
		RateLimitIP:    cfg.RateLimitIP,
		RateLimitToken: cfg.RateLimitToken,
		BlockDuration:  cfg.BlockDuration,
	})

	h := handler.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.Health)
	mux.HandleFunc("/api/v1/example", h.Example)

	// Aplica o middleware de rate limit sobre o mux
	handlerWithMiddleware := ratemw.RateLimiter(l)(mux)

	addr := fmt.Sprintf(":%s", cfg.HTTPPort)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatal(err)
	}
}