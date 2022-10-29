package request_handler

import (
	partition "dbcache/partitioning"
	"net/rpc"
	"fmt"
)

type Port uint16

func (p Port) String() string {
	return fmt.Sprintf("%d", uint32(p))
}

type PeerInfo struct {
	Name   string
	Port   Port
}

type Peer struct {
	info    PeerInfo
	putChan chan putReq
	reqChan chan CacheRequest
	conn   *rpc.Client
}

// type ClusterNodes map[PeerInfo]*Peer

// type ClusterNodesDeque deque.Deque

type Cluster struct {
	info   map[Peer]bool
	seeder Peer
	nodes map[partition.Partition]*ClusterNodes
	sortedPartitions []partition.Partition
}

type ShareInfoResponse struct {
	Info       map[Peer]bool
	Partitions map[partition.Partition]map[PeerInfo]bool
}

type ShareCacheRequest struct{}

type CacheRequest struct {
	Action int8
	Key string
	Value string
}

type GetRequest string

type CacheRequestResponse struct {
	Ok bool
	Key string
	Value string
}

type putReq struct {
	key string
	val string
}