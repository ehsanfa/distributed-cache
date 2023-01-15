package buddy

import (
	"dbcache/cluster/peer"
)

type Buddies interface {
	Add(peer.Peer) bool
	All() map[peer.Peer]bool
	Count() int
	IsEmpty() bool
	IsBuddyWith(peer.Peer) bool
	CanAcceptBuddyRequest() bool
	Remove(peer.Peer)
}
