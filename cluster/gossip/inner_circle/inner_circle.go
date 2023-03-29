package innercircle

import (
	"dbcache/cluster/peer"
)

type InnerCircle interface {
	All() map[peer.Peer]bool
	Add(peer.Peer) bool
	Count() int
	IsEmpty() bool
	Remove(peer.Peer)
}
