package limiter

import (
	"fmt"
	"time"
)

type Config struct {
	RateLimitIP    int
	RateLimitToken int
	BlockDuration  time.Duration
}

type Limiter struct {
	store Store
	cfg   Config
}

func NewLimiter(store Store, cfg Config) *Limiter {
	return &Limiter{store: store, cfg: cfg}
}

type Result struct {
	Allowed bool
	Reason  string
}

func (l *Limiter) Allow(ip, token string) (*Result, error) {
	var keyType string
	var limit int

	if token != "" {
		keyType = "token"
		limit = l.cfg.RateLimitToken
	} else {
		keyType = "ip"
		limit = l.cfg.RateLimitIP
	}

	identifier := token
	if identifier == "" {
		identifier = ip
	}

	if identifier == "" {
		return &Result{Allowed: true}, nil
	}

	blockKey := l.blockKey(keyType, identifier)
	rateKey := l.rateKey(keyType, identifier)

	blocked, err := l.store.IsBlocked(blockKey)
	if err != nil {
		return nil, err
	}
	if blocked {
		return &Result{Allowed: false, Reason: "blocked"}, nil
	}

	count, err := l.store.Increment(rateKey, 1)
	if err != nil {
		return nil, err
	}

	if int(count) > limit {
		if err := l.store.Block(blockKey, l.cfg.BlockDuration); err != nil {
			return nil, err
		}
		return &Result{Allowed: false, Reason: "limit-exceeded"}, nil
	}

	return &Result{Allowed: true}, nil
}

func (l *Limiter) blockKey(t, id string) string {
	return fmt.Sprintf("block:%s:%s", t, id)
}

func (l *Limiter) rateKey(t, id string) string {
	return fmt.Sprintf("rate:%s:%s", t, id)
}
