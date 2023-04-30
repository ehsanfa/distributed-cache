package network

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
)

type Node interface {
	Introduce(peer.Peer) error
	GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error)
	GetCache() (map[string]cacher.CacheValue, error)
	AskForParition() (partition.Partition, error)
	UpdateBuffer(buffer.Buffer) error
	Ping() (bool, error)
	Get(string) (cacher.CacheValue, error)
	Set(key string, value cacher.CacheValue) error
}
