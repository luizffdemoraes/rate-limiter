package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad_UsesDefaultsWhenEnvNotSet(t *testing.T) {
	// Garante que não há variáveis setadas
	_ = os.Unsetenv("HTTP_PORT")
	_ = os.Unsetenv("RATE_LIMIT_IP")
	_ = os.Unsetenv("RATE_LIMIT_TOKEN")
	_ = os.Unsetenv("BLOCK_DURATION")
	_ = os.Unsetenv("REDIS_HOST")
	_ = os.Unsetenv("REDIS_PORT")
	_ = os.Unsetenv("REDIS_DB")
	_ = os.Unsetenv("REDIS_PASSWORD")

	cfg := Load()

	require.Equal(t, "8080", cfg.HTTPPort)
	require.Equal(t, 10, cfg.RateLimitIP)
	require.Equal(t, 100, cfg.RateLimitToken)
	require.Equal(t, 300*time.Second, cfg.BlockDuration)
	require.Equal(t, "redis", cfg.RedisHost)
	require.Equal(t, "6379", cfg.RedisPort)
	require.Equal(t, 0, cfg.RedisDB)
	require.Equal(t, "", cfg.RedisPass)
}

func TestLoad_UsesEnvValuesWhenSet(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("RATE_LIMIT_IP", "3")
	t.Setenv("RATE_LIMIT_TOKEN", "7")
	t.Setenv("BLOCK_DURATION", "60")
	t.Setenv("REDIS_HOST", "custom-redis")
	t.Setenv("REDIS_PORT", "6380")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("REDIS_PASSWORD", "secret")

	cfg := Load()

	require.Equal(t, "9090", cfg.HTTPPort)
	require.Equal(t, 3, cfg.RateLimitIP)
	require.Equal(t, 7, cfg.RateLimitToken)
	require.Equal(t, 60*time.Second, cfg.BlockDuration)
	require.Equal(t, "custom-redis", cfg.RedisHost)
	require.Equal(t, "6380", cfg.RedisPort)
	require.Equal(t, 2, cfg.RedisDB)
	require.Equal(t, "secret", cfg.RedisPass)
}

