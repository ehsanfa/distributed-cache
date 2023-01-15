package node

import (
	"dbcache/cluster/buddy"
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/connection"
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"net/rpc"
)

const maxBuddyNum = 2

type Node struct {
	bufferSizeExceeded chan bool
	cacheVersions      map[string]cacher.CacheVersion
	connections        map[peer.Peer]*rpc.Client
	connection         connection.Connection
	partitions         []partition.Partition
	partition          partition.Partition
	isSeeder           bool
	buddies            buddy.Buddies
	seeder             peer.Peer
	info               info.ClusterInfo
	buffer             buffer.Buffer
	cache              cacher.Cache
	peer               peer.Peer
}

func CreateNode(p peer.Peer, isSeeder bool, partition partition.Partition) *Node {
	return &Node{
		peer:        p,
		isSeeder:    isSeeder,
		partition:   partition,
		info:        info.CreateInMemoryClusterInfo(),
		cache:       cacher.CreateInMemoryCache(),
		buffer:      buffer.CreateInMemoryBuffer(),
		connections: make(map[peer.Peer]*rpc.Client),
		buddies:     buddy.CreateInMemoryBuddies(maxBuddyNum),
	}
}

func (n *Node) SetSeeder(p peer.Peer) {
	n.seeder = p
}
