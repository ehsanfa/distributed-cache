package innercircle

import (
	"dbcache/cluster/peer"
	"sync"

	"github.com/gammazero/deque"
)

type InMemoryBuddies struct {
	buddies     map[peer.Peer]bool
	maxBuddyNum int
	mu          sync.RWMutex
	dq          *deque.Deque[peer.Peer]
}

func CreateInMemoryBuddies(maxBuddyNum int) *InMemoryBuddies {
	return &InMemoryBuddies{
		buddies:     make(map[peer.Peer]bool),
		maxBuddyNum: maxBuddyNum,
		dq:          deque.New[peer.Peer](maxBuddyNum),
	}
}

func (b *InMemoryBuddies) replaceDeque() {
	b.dq = deque.New[peer.Peer](b.maxBuddyNum)
	for p := range b.All() {
		b.dq.PushFront(p)
	}
}

func (b *InMemoryBuddies) canAdd(p peer.Peer) bool {
	return b.canAcceptBuddyRequest() && !b.isBuddyWith(p)
}

func (b *InMemoryBuddies) Add(p peer.Peer) bool {
	if !b.canAdd(p) {
		return false
	}
	b.mu.Lock()
	b.buddies[p] = true
	b.replaceDeque()
	b.mu.Unlock()
	return true
}

func (b *InMemoryBuddies) All() map[peer.Peer]bool {
	return b.buddies
}

func (b *InMemoryBuddies) Count() int {
	buddies := b.All()
	return len(buddies)
}

func (b *InMemoryBuddies) IsEmpty() bool {
	return b.Count() == 0
}

func (b *InMemoryBuddies) isBuddyWith(p peer.Peer) bool {
	buddies := b.All()
	_, ok := buddies[p]
	return ok
}

func (b *InMemoryBuddies) canAcceptBuddyRequest() bool {
	return b.Count() < b.maxBuddyNum
}

func (b *InMemoryBuddies) Remove(p peer.Peer) {
	b.mu.Lock()
	delete(b.buddies, p)
	b.replaceDeque()
	b.mu.Unlock()
}

func (b *InMemoryBuddies) Replace(old peer.Peer, new peer.Peer) {
	b.Remove(old)
	b.Add(new)
}

func (b *InMemoryBuddies) Diff(target map[peer.Peer]bool) map[peer.Peer]bool {
	diff := make(map[peer.Peer]bool)
	for p := range target {
		if !b.isBuddyWith(p) {
			diff[p] = true
		}
	}
	return diff
}

func (b *InMemoryBuddies) Shuffle(candidates map[peer.Peer]bool) {
	// find the difference between candidates and the current circle
	diff := b.Diff(candidates)
	// if the difference is empty, do nothing
	if len(diff) == 0 {
		return
	}
	// if the count is
	counter := 0
	dq := *b.dq
	for newP := range diff {
		if b.isBuddyWith(newP) {
			continue
		}
		if counter == b.maxBuddyNum || dq.Len() == 0 {
			break
		}
		oldP := dq.PopBack()
		if b.canAdd(newP) {
			b.Add(newP)
		} else {
			b.Replace(oldP, newP)
		}
		counter++
	}
}
