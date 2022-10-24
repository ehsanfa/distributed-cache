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

type Peer struct {
	Name   string
	Port   Port
	putChan chan putReq
	reqChan chan CacheRequest
	conn   *rpc.Client
}

type Cluster struct {
	info   map[Peer]bool
	seeder Peer
	partitions map[partition.Partition]map[*Peer]*rpc.Client
	sortedPartitions []partition.Partition
}

type ShareInfoResponse struct {
	Info       map[Peer]bool
	Partitions map[partition.Partition]map[Peer]bool
}

type ShareCacheRequest struct{}

type CacheRequest struct {
	Action int8
	Key string
	Value string
}

type CacheRequestResponse struct {
	Ok bool
	Key string
	Value string
}

type putReq struct {
	key string
	val string
}