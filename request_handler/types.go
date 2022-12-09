package request_handler

import (
	partition "dbcache/partitioning"
	"fmt"
	"net/rpc"
	"sync"
)

type Port uint16

func (p Port) String() string {
	return fmt.Sprintf("%d", uint32(p))
}

type PeerInfo struct {
	Name string
	Port Port
}

type Peer struct {
	info    PeerInfo
	putChan chan putReq
	reqChan chan CacheRequest
	conn    *rpc.Client
}

// type ClusterNodes map[PeerInfo]*Peer

// type ClusterNodesDeque deque.Deque

type Cluster struct {
	sync.RWMutex
	info             map[Peer]bool
	seeder           Peer
	nodes            map[partition.Partition]*ClusterNodes
	sortedPartitions []partition.Partition
}

type ShareInfoResponse struct {
	Info       map[Peer]bool
	Partitions map[partition.Partition]map[PeerInfo]bool
}

type ShareCacheRequest struct{}

type CacheRequest struct {
	Action int8
	Key    string
	Value  string
}

type GetRequest string

type GetCacheResponse struct {
	Ok    bool
	Value string
}

type PutCacheResponse bool

type putReq struct {
	key string
	val string
}
