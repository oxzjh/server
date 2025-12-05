package rate

import (
	"sync"
	"time"
)

type Group struct {
	sync.RWMutex
	limit    time.Duration
	burst    int
	limiters map[string]*Limiter
}

func (g *Group) Allow(key string) bool {
	g.Lock()
	lim, ok := g.limiters[key]
	if !ok {
		lim = NewLimiter(g.limit, g.burst)
		g.limiters[key] = lim
	}
	g.Unlock()
	return lim.Allow()
}

func (g *Group) GetLimit(key string) time.Duration {
	g.RLock()
	defer g.RUnlock()
	if lim, ok := g.limiters[key]; ok {
		return lim.limit
	}
	return g.limit
}

func (g *Group) SetLimit(key string, limit time.Duration) {
	g.RLock()
	defer g.RUnlock()
	if lim, ok := g.limiters[key]; ok {
		lim.limit = limit
	}
}

func NewGroup(limit time.Duration, burst int) *Group {
	return &Group{limit: limit, burst: burst, limiters: make(map[string]*Limiter)}
}
