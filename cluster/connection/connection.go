package connection

import (
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
)

type Connection interface {
	Introduce(peer.Peer) (IntroductionResponse, error)
}

type IntroductionResponse interface {
	ClusterInfo() info.ClusterInfo
	Cache() cacher.Cache
	Partition() partition.Partition
}
