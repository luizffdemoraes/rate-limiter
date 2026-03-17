package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/lffm1994/rate-limiter/internal/limiter"
)


type inMemoryStoreMiddleware struct {
	counters map[string]int64
	blocked  map[string]time.Time
	now      func() time.Time
}

func newInMemoryStoreMiddleware() *inMemoryStoreMiddleware {
	return &inMemoryStoreMiddleware{
		counters: make(map[string]int64),
		blocked:  make(map[string]time.Time),
		now:      time.Now,
	}
}

func (s *inMemoryStoreMiddleware) Increment(key string, _ int) (int64, error) {
	s.counters[key]++
	return s.counters[key], nil
}

func (s *inMemoryStoreMiddleware) IsBlocked(key string) (bool, error) {
	until, ok := s.blocked[key]
	if !ok {
		return false, nil
	}
	if s.now().After(until) {
		delete(s.blocked, key)
		return false, nil
	}
	return true, nil
}

func (s *inMemoryStoreMiddleware) Block(key string, duration time.Duration) error {
	s.blocked[key] = s.now().Add(duration)
	return nil
}

func TestRateLimiterMiddleware_Returns429AndBodyOnLimitExceeded(t *testing.T) {
	store := newInMemoryStoreMiddleware()
	l := limiter.NewLimiter(store, limiter.Config{
		RateLimitIP:    1,
		RateLimitToken: 1,
		BlockDuration:  time.Minute,
	})

	called := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	})

	mw := RateLimiter(l)(next)

	// Primeira requisição: permitida
	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, 1, called)

	// Segunda requisição no mesmo IP: deve estourar limite e retornar 429 com corpo exato
	req2 := httptest.NewRequest(http.MethodGet, "/any", nil)
	req2.RemoteAddr = "127.0.0.1:12345"
	rr2 := httptest.NewRecorder()
	mw.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusTooManyRequests, rr2.Code)
	require.Equal(t, tooManyRequestsBody, rr2.Body.String())
	// Não deve chamar o next handler novamente
	require.Equal(t, 1, called)
}

