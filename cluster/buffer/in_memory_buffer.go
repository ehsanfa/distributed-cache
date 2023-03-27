package buffer

import (
	"dbcache/cluster/cacher"
	"sync"

	ll "github.com/ehsanfa/linked-list"
)

type InMemoryBuffer struct {
	internal      ll.LinkedList
	sharingBuffer map[string]cacher.CacheValue
	mu            sync.RWMutex
}

type cachable struct {
	key   string
	value cacher.CacheValue
}

func (c *cachable) GetKey() string {
	return c.key
}

func (c *cachable) GetValue() cacher.CacheValue {
	return c.value
}

func CreateInMemoryBuffer() *InMemoryBuffer {
	return &InMemoryBuffer{}
}

func (b *InMemoryBuffer) IsEmpty() bool {
	b.mu.RLock()
	isEmpty := b.internal.Count() == 0 && len(b.sharingBuffer) == 0
	b.mu.RUnlock()
	return isEmpty
}

func (b *InMemoryBuffer) Add(c cacher.Cachable) {
	b.mu.Lock()
	b.internal.Append(c)
	if b.sharingBuffer == nil {
		b.sharingBuffer = make(map[string]cacher.CacheValue)
	}
	b.sharingBuffer[c.GetKey()] = c.GetValue()
	b.mu.Unlock()
}

func (b *InMemoryBuffer) All() map[string]cacher.CacheValue {
	b.mu.RLock()
	s := b.sharingBuffer
	b.mu.RUnlock()
	return s
}

func (b *InMemoryBuffer) Reset() {
	e := CreateInMemoryBuffer()
	b.mu.Lock()
	b.internal = e.internal
	b.sharingBuffer = e.sharingBuffer
	b.mu.Unlock()
}

func (b *InMemoryBuffer) Size() int {
	b.mu.RLock()
	c := b.internal.Count()
	b.mu.RUnlock()
	return c
}

func (b *InMemoryBuffer) Merge(source Buffer) {
	for k, v := range source.All() {
		b.Add(&cachable{k, v})
	}
}
