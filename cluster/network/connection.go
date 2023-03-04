package network

import (
	"dbcache/cluster/cacher"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
)

type Server interface {
	GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error)
	GetCache() (map[string]cacher.CacheValue, error)
	AskForParition() (partition.Partition, error)
	Ping() (bool, error)
	peer.WithPort
}
