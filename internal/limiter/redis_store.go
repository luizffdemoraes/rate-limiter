package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(host, port, pass string, db int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       db,
	})
	return &RedisStore{client: client}
}

func (s *RedisStore) Increment(key string, windowSeconds int) (int64, error) {
	ctx := context.Background()

	val, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if val == 1 {
		_ = s.client.Expire(ctx, key, time.Duration(windowSeconds)*time.Second).Err()
	}
	return val, nil
}

func (s *RedisStore) IsBlocked(key string) (bool, error) {
	ctx := context.Background()
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (s *RedisStore) Block(key string, duration time.Duration) error {
	ctx := context.Background()
	return s.client.Set(ctx, key, "1", duration).Err()
}
