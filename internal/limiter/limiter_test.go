package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type inMemoryStore struct {
	counters map[string]int64
	blocked  map[string]time.Time
	now      func() time.Time
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		counters: make(map[string]int64),
		blocked:  make(map[string]time.Time),
		now:      time.Now,
	}
}

func (s *inMemoryStore) Increment(key string, _ int) (int64, error) {
	s.counters[key]++
	return s.counters[key], nil
}

func (s *inMemoryStore) IsBlocked(key string) (bool, error) {
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

func (s *inMemoryStore) Block(key string, duration time.Duration) error {
	s.blocked[key] = s.now().Add(duration)
	return nil
}

func TestLimiter_IPLimitAndBlock(t *testing.T) {
	store := newInMemoryStore()
	l := NewLimiter(store, Config{
		RateLimitIP:    2,
		RateLimitToken: 5,
		BlockDuration:  time.Minute,
	})

	ip := "1.1.1.1"

	// Primeiras 2 requisições devem ser permitidas
	for i := 0; i < 2; i++ {
		res, err := l.Allow(ip, "")
		require.NoError(t, err)
		require.True(t, res.Allowed, "request %d should be allowed", i+1)
	}

	// Terceira no mesmo segundo deve exceder limite → bloqueio
	res, err := l.Allow(ip, "")
	require.NoError(t, err)
	require.False(t, res.Allowed)
	require.Equal(t, "limit-exceeded", res.Reason)

	// Após bloqueio, novas requisições continuam negadas
	res, err = l.Allow(ip, "")
	require.NoError(t, err)
	require.False(t, res.Allowed)
	require.Equal(t, "blocked", res.Reason)
}

func TestLimiter_TokenHasPrecedenceOverIP(t *testing.T) {
	store := newInMemoryStore()
	l := NewLimiter(store, Config{
		RateLimitIP:    2,
		RateLimitToken: 4,
		BlockDuration:  time.Minute,
	})

	ip := "2.2.2.2"
	token := "my-token"

	// Mesmo IP, mas com token: deve respeitar limite do token (4), não o de IP (2)
	for i := 0; i < 4; i++ {
		res, err := l.Allow(ip, token)
		require.NoError(t, err)
		require.True(t, res.Allowed, "token request %d should be allowed", i+1)
	}

	// Quinta requisição com o mesmo token deve estourar limite do token
	res, err := l.Allow(ip, token)
	require.NoError(t, err)
	require.False(t, res.Allowed)
	require.Equal(t, "limit-exceeded", res.Reason)
}

