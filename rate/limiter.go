package rate

import "time"

type Limiter struct {
	limit  time.Duration
	burst  int
	tokens int
	last   time.Time
}

func (lim *Limiter) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(lim.last)
	if elapsed > lim.limit {
		delta := elapsed / lim.limit
		lim.tokens += int(delta)
		if lim.tokens > lim.burst {
			lim.tokens = lim.burst
			lim.last = now
		} else {
			lim.last = lim.last.Add(delta * lim.limit)
		}
	}
	if lim.tokens > 0 {
		lim.tokens--
		return true
	}
	return false
}

func NewLimiter(limit time.Duration, burst int) *Limiter {
	return &Limiter{limit: limit, burst: burst}
}
