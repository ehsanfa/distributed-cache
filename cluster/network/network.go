package network

import (
	"dbcache/cluster/peer"
	"time"
)

type Network interface {
	Connect(peer peer.Peer, timeout time.Duration) (Node, error)
	Peer() peer.Peer
	Kill()
}
