package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Option func(*redis.Pool)

func WithMaxIdle(maxIdle int) Option {
	return func(o *redis.Pool) {
		o.MaxIdle = maxIdle
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *redis.Pool) {
		o.IdleTimeout = idleTimeout
	}
}
