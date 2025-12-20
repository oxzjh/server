package cache

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type ICache interface {
	Set(string, any, int64)
	Get(string) string
	GetInt64(string) int64
	GetFloat64(string) float64
	Incr(string) int64
	Incrby(string, int64) int64
	Decr(string) int64
	Decrby(string, int64) int64
	Del(string)
	Hset(string, string, string)
	Hget(string, string) string
	Hdel(string, string)
	Hgetall(string) map[string]string
	Hmset(string, string, string, string, string, ...any)
	Save()
	Close()
}

type cache struct {
	sync.Mutex
	file string
	Ex   map[string]int64             `json:"e"`
	Data map[string]any               `json:"d"`
	Hash map[string]map[string]string `json:"h"`
}

func (c *cache) Set(key string, val any, expire int64) {
	c.Lock()
	defer c.Unlock()
	switch v := val.(type) {
	case int64, string, float64:
	case int8:
		val = int64(v)
	case uint8:
		val = int64(v)
	case int16:
		val = int64(v)
	case uint16:
		val = int64(v)
	case int:
		val = int64(v)
	case uint:
		val = int64(v)
	case uint64:
		val = int64(v)
	case float32:
		val = float64(v)
	default:
		log.Println("unsupport type:", key)
		return
	}
	c.Data[key] = val
	if expire > 0 {
		c.Ex[key] = time.Now().Unix() + expire
	}
}

func (c *cache) Get(key string) string {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.Data[key]; ok && !c.expired(key) {
		switch v := val.(type) {
		case string:
			return v
		default:
			log.Println("type error:", key)
		}
	}
	return ""
}

func (c *cache) GetInt64(key string) int64 {
	c.Lock()
	defer c.Unlock()
	return c.getInt(key)
}

func (c *cache) GetFloat64(key string) float64 {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.Data[key]; ok && !c.expired(key) {
		switch v := val.(type) {
		case float64:
			return v
		default:
			log.Println("type error:", key)
		}
	}
	return 0
}

func (c *cache) Incr(key string) int64 {
	c.Lock()
	defer c.Unlock()
	val := c.getInt(key)
	val++
	c.Data[key] = val
	return val
}

func (c *cache) Incrby(key string, increment int64) int64 {
	c.Lock()
	defer c.Unlock()
	val := c.getInt(key)
	val += increment
	c.Data[key] = val
	return val
}

func (c *cache) Decr(key string) int64 {
	c.Lock()
	defer c.Unlock()
	val := c.getInt(key)
	val--
	c.Data[key] = val
	return val
}

func (c *cache) Decrby(key string, decrement int64) int64 {
	c.Lock()
	defer c.Unlock()
	val := c.getInt(key)
	val -= decrement
	c.Data[key] = val
	return val
}

func (c *cache) Del(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Ex, key)
	delete(c.Data, key)
}

func (c *cache) Hset(key, sub, val string) {
	c.Lock()
	defer c.Unlock()
	m, ok := c.Hash[key]
	if !ok {
		m = make(map[string]string)
		c.Hash[key] = m
	}
	m[sub] = val
}

func (c *cache) Hget(key, sub string) string {
	c.Lock()
	defer c.Unlock()
	if m, ok := c.Hash[key]; ok {
		return m[sub]
	}
	return ""
}

func (c *cache) Hdel(key, sub string) {
	c.Lock()
	defer c.Unlock()
	if m, ok := c.Hash[key]; ok {
		delete(m, sub)
	}
}

func (c *cache) Hgetall(key string) map[string]string {
	c.Lock()
	defer c.Unlock()
	return c.Hash[key]
}

func (c *cache) Hmset(key, sub1, val1, sub2, val2 string, args ...any) {
	c.Lock()
	defer c.Unlock()
	val := map[string]string{sub1: val1, sub2: val2}
	l := len(args)
	for i := 0; i < l; i += 2 {
		val[args[i].(string)] = args[i+1].(string)
	}
	c.Hash[key] = val
}

func (c *cache) Save() {
	b, _ := json.Marshal(c)
	os.WriteFile(c.file, b, 0644)
}

func (*cache) Close() {
}

func (c *cache) expired(key string) bool {
	if expire, ok := c.Ex[key]; ok && time.Now().Unix() > expire {
		delete(c.Ex, key)
		delete(c.Data, key)
		return true
	}
	return false
}

func (c *cache) getInt(key string) int64 {
	if val, ok := c.Data[key]; ok && !c.expired(key) {
		switch v := val.(type) {
		case int64:
			return v
		case float64:
			return int64(v)
		default:
			log.Println("type error:", key)
		}
	}
	return 0
}

func New(file string) ICache {
	c := &cache{file: file}
	f, err := os.Open(file)
	if err != nil {
		c.Ex = make(map[string]int64)
		c.Data = make(map[string]any)
		c.Hash = make(map[string]map[string]string)
	} else {
		json.NewDecoder(f).Decode(c)
		f.Close()
		if c.Ex == nil {
			c.Ex = make(map[string]int64)
		}
		if c.Data == nil {
			c.Data = make(map[string]any)
		}
		if c.Hash == nil {
			c.Hash = make(map[string]map[string]string)
		}
	}
	return c
}
