package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort       string
	RateLimitIP    int
	RateLimitToken int
	BlockDuration  time.Duration

	RedisHost string
	RedisPort string
	RedisDB   int
	RedisPass string
}

func Load() *Config {
	_ = godotenv.Overload()

	return &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
		RateLimitIP:    getEnvInt("RATE_LIMIT_IP", 10),
		RateLimitToken: getEnvInt("RATE_LIMIT_TOKEN", 100),
		BlockDuration:  time.Duration(getEnvInt("BLOCK_DURATION", 300)) * time.Second,
		RedisHost:      getEnv("REDIS_HOST", "redis"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		RedisPass:      getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("invalid value for %s (%s), using default %d", key, v, def)
			return def
		}
		return i
	}
	return def
}
