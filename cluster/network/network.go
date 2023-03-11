package network

import (
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
)

type Network interface {
	Connect(peer.Peer) Node
	Serve(p peer.Peer, info info.ClusterInfoProvider) (peer.WithPort, error)
	Kill()
}

type Node interface {
	GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error)
	GetCache() (map[string]cacher.CacheValue, error)
	AskForParition() (partition.Partition, error)
	Ping() (bool, error)
	// peer.WithPort
}
