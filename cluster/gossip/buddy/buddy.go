package buddy

import (
	"dbcache/cluster/peer"
)

type Buddy peer.Peer

type HasBuddies interface {
	AllBuddies() map[Buddy]bool
}

type Buddies interface {
	HasBuddies
	Add(peer.Peer) bool
	Count() int
	IsEmpty() bool
	IsBuddyWith(peer.Peer) bool
	CanAcceptBuddyRequest() bool
	Remove(peer.Peer)
}
