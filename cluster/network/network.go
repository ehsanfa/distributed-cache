package network

import (
	"dbcache/cluster/info"
	"dbcache/cluster/peer"
)

type Network interface {
	Connect(peer.Peer) Node
	Serve(p peer.Peer, info info.ClusterInfoProvider) (peer.WithPort, error)
	Kill()
}
