package cache

import (
	"sync"
	"time"
)

type wrapper struct {
	value     any
	createdAt time.Time
}

type Cache struct {
	mu   sync.Mutex
	impl map[string]*wrapper
	ttl  time.Duration
}

const cleanCacheFrequency = 100

func New(TTL time.Duration) *Cache {
	return &Cache{impl: map[string]*wrapper{}, ttl: TTL}
}

func (c *Cache) conditionClean() {
	c.ttl++
	if c.ttl < cleanCacheFrequency {
		return
	}

	c.ttl = 0
	toRemove := make([]string, 0, len(c.impl))
	now := time.Now()

	for k, v := range c.impl {
		if now.Sub(v.createdAt) > c.ttl {
			toRemove = append(toRemove, k)
		}
	}

	for _, k := range toRemove {
		delete(c.impl, k)
	}
}

func (c *Cache) Store(id string, v any) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conditionClean()
	c.impl[id] = &wrapper{value: v, createdAt: now}
}

func (c *Cache) Load(id string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conditionClean()

	v, ok := c.impl[id]
	if ok {
		return v.value, true
	}
	return nil, false
}
