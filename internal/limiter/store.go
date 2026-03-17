package limiter

import "time"

type Store interface {
	Increment(key string, windowSeconds int) (int64, error)
	IsBlocked(key string) (bool, error)
	Block(key string, duration time.Duration) error
}
