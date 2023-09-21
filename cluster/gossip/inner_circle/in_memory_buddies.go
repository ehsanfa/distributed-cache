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
	self        peer.Peer
}

func CreateInMemoryBuddies(maxBuddyNum int, self peer.Peer) *InMemoryBuddies {
	return &InMemoryBuddies{
		buddies:     make(map[peer.Peer]bool),
		maxBuddyNum: maxBuddyNum,
		dq:          deque.New[peer.Peer](maxBuddyNum),
		self:        self,
	}
}

func (b *InMemoryBuddies) replaceDeque() {
	b.dq = deque.New[peer.Peer](b.maxBuddyNum)
	for p := range b.All() {
		b.dq.PushFront(p)
	}
}

func (b *InMemoryBuddies) canAdd(p peer.Peer) bool {
	return b.canAcceptBuddyRequest() && !b.isBuddyWith(p) && p != nil && p != b.self
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

func (b *InMemoryBuddies) Diff(target []peer.Peer) []peer.Peer {
	diff := make([]peer.Peer, 0)
	for _, p := range target {
		if !b.isBuddyWith(p) {
			diff = append(diff, p)
		}
	}
	return diff
}

func (b *InMemoryBuddies) Shuffle(candidates []peer.Peer) {
	diff := b.Diff(candidates)
	if len(diff) == 0 {
		return
	}
	if len(b.buddies) == 0 {
		b.initiateBuddiesFromCandidates(candidates)
		return
	}
	counter := 0
	dq := *b.dq
	for _, newP := range diff {
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

func (b *InMemoryBuddies) initiateBuddiesFromCandidates(candidates []peer.Peer) {
	for _, p := range candidates {
		if b.canAdd(p) {
			b.Add(p)
		}
	}
}
