package network

import "dbcache/cluster/peer"

type Cluster interface {
	Connect(peer.Peer) Server
}
