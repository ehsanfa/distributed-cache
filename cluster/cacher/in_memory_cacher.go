package cacher

import "sync"

type InMemoryCache struct {
	cache map[string]CacheValue
	mu    sync.RWMutex
}

func (c *InMemoryCache) Get(key string) (value CacheValue, ok bool) {
	c.mu.RLock()
	value, ok = c.cache[key]
	c.mu.RUnlock()
	return
}

func (c *InMemoryCache) Exists(key string) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *InMemoryCache) Set(key string, value CacheValue) error {
	c.mu.Lock()
	c.cache[key] = value
	c.mu.Unlock()
	return nil
}

func (c *InMemoryCache) All() map[string]CacheValue {
	return c.cache
}

func (c *InMemoryCache) Touch(key string) {
	v, ok := c.Get(key)
	if !ok {
		v = NewVersionBasedCacheValue("", 1)
	}
	v.IncrementVersion()
	c.Set(key, v)
}

func (c *InMemoryCache) Version(key string) int {
	v, ok := c.Get(key)
	if !ok {
		return 0
	}
	return v.Version()
}

func (c *InMemoryCache) Replace(cache map[string]CacheValue) {
	c.mu.Lock()
	c.cache = cache
	c.mu.Unlock()
}

func (c *InMemoryCache) Delete(key string) {
	v, ok := c.Get(key)
	if !ok {
		return
	}
	v.SetValue("")
	v.IncrementVersion()
	c.Set(key, v)
}

func CreateInMemoryCache() *InMemoryCache {
	return &InMemoryCache{cache: make(map[string]CacheValue)}
}
