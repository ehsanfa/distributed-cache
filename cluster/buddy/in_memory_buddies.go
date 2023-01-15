package buddy

import (
	"dbcache/cluster/peer"
	"sync"
)

type InMemoryBuddies struct {
	buddies     map[peer.Peer]bool
	maxBuddyNum int
	mu          sync.RWMutex
}

func CreateInMemoryBuddies(maxBuddyNum int) *InMemoryBuddies {
	return &InMemoryBuddies{
		buddies:     make(map[peer.Peer]bool),
		maxBuddyNum: maxBuddyNum,
	}
}

func (b *InMemoryBuddies) Add(p peer.Peer) bool {
	if !b.CanAcceptBuddyRequest() {
		return false
	}
	b.mu.Lock()
	b.buddies[p] = true
	b.mu.Unlock()
	return true
}

func (b *InMemoryBuddies) All() map[peer.Peer]bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buddies
}

func (b *InMemoryBuddies) Count() int {
	buddies := b.All()
	return len(buddies)
}

func (b *InMemoryBuddies) IsEmpty() bool {
	return b.Count() == 0
}

func (b *InMemoryBuddies) IsBuddyWith(p peer.Peer) bool {
	buddies := b.All()
	_, ok := buddies[p]
	return ok
}

func (b *InMemoryBuddies) CanAcceptBuddyRequest() bool {
	return b.Count() < b.maxBuddyNum
}

func (b *InMemoryBuddies) Remove(p peer.Peer) {
	b.mu.Lock()
	delete(b.buddies, p)
	b.mu.Unlock()
}
