package node

// import (
// 	"dbcache/cluster/buffer"
// 	"dbcache/cluster/cacher"
// 	"dbcache/cluster/gossip/buddy"
// 	"dbcache/cluster/gossip/gossip"
// 	"dbcache/cluster/info"
// 	"dbcache/cluster/network"
// 	"dbcache/cluster/partition"
// 	"dbcache/cluster/peer"
// 	"net/rpc"
// )

// const maxBuddyNum = 2

// type Node struct {
// 	bufferSizeExceeded chan bool
// 	// cacheVersions      map[string]cacher.CacheVersion
// 	connections map[peer.Peer]*rpc.Client
// 	partitions  []partition.Partition
// 	// partition          partition.Partition
// 	network  network.Network
// 	isSeeder bool
// 	buddies  buddy.Buddies
// 	gossip   gossip.Gossip
// 	seeder   peer.Peer
// 	info     info.ClusterInfo
// 	buffer   buffer.Buffer
// 	cache    cacher.Cache
// 	peer     peer.Peer
// }

// func (n *Node) Introduce() error {
// 	var err error
// 	// Get Info
// 	i, err := n.network.Connect(n.seeder).GetClusterInfo()
// 	if err != nil {
// 		return err
// 	}
// 	n.info.Replace(i)

// 	// Get Cache
// 	c, err := n.network.Connect(n.seeder).GetCache()
// 	if err != nil {
// 		return err
// 	}
// 	n.cache.Replace(c)

// 	// Ask for a partition
// 	suggestedPartition, err := n.network.Connect(n.seeder).AskForParition()
// 	if err != nil {
// 		return err
// 	}
// 	n.peer.SetPartition(suggestedPartition)
// 	return nil
// }

// func CreateNode(p peer.Peer, isSeeder bool) *Node {
// 	return &Node{
// 		peer:     p,
// 		isSeeder: isSeeder,
// 		// partition:   partition,
// 		info:        info.CreateInMemoryClusterInfo(),
// 		cache:       cacher.CreateInMemoryCache(),
// 		buffer:      buffer.CreateInMemoryBuffer(),
// 		connections: make(map[peer.Peer]*rpc.Client),
// 		buddies:     buddy.CreateInMemoryBuddies(maxBuddyNum),
// 	}
// }

// func (n *Node) SetSeeder(p peer.Peer) {
// 	n.seeder = p
// }
