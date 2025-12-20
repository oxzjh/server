package redis

import (
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/oxzjh/server/cache"
)

type redisCache struct {
	pool         *redis.Pool
	errorHandler func(error, string)
}

func (r *redisCache) Set(key string, val any, expire int64) {
	conn := r.pool.Get()
	defer conn.Close()
	if expire > 0 {
		conn.Do("SET", key, val, "EX", expire)
	} else {
		conn.Do("SET", key, val)
	}
}

func (r *redisCache) Get(key string) string {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) GetInt64(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("GET", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) GetFloat64(key string) float64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Float64(conn.Do("GET", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Incr(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("INCR", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Incrby(key string, increment int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("INCRBY", key, increment))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Decr(key string) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("DECR", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Decrby(key string, decrement int64) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("DECRBY", key, decrement))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Del(key string) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("DEL", key)
}

func (r *redisCache) Hset(key, sub, val string) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("HSET", key, sub, val)
}

func (r *redisCache) Hget(key, sub string) string {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("HGET", key, sub))
	r.errorHandler(err, key+"."+sub)
	return val
}

func (r *redisCache) HgetInt64(key, sub string) int64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("HGET", key, sub))
	r.errorHandler(err, key+"."+sub)
	return val
}

func (r *redisCache) HgetFloat64(key, sub string) float64 {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Float64(conn.Do("HGET", key, sub))
	r.errorHandler(err, key+"."+sub)
	return val
}

func (r *redisCache) Hdel(key, sub string) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("HDEL", key, sub)
}

func (r *redisCache) Hgetall(key string) map[string]string {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.StringMap(conn.Do("HGETALL", key))
	r.errorHandler(err, key)
	return val
}

func (r *redisCache) Hmset(key, sub1, val1, sub2, val2 string, args ...any) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Do("HMSET", append([]any{key, sub1, val1, sub2, val2}, args...)...)
}

func (*redisCache) Save() {
}

func (r *redisCache) Close() {
	r.pool.Close()
}

func New(addr string, db uint8, opts ...Option) cache.ICache {
	pool := &redis.Pool{
		MaxIdle: 2,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if _, err := conn.Do("SELECT", db); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
	}
	for _, opt := range opts {
		opt(pool)
	}
	return &redisCache{pool, func(err error, key string) {
		log.Println("redis:", key, err)
	}}
}
