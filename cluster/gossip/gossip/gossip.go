package gossip

import (
	"dbcache/cluster/cluster"
	"dbcache/cluster/gossip/buddy"
	"dbcache/cluster/info"
	"dbcache/cluster/peer"
)

type Gossip interface {
	Initialize(cluster.Cluster, info.ClusterInfo, buddy.HasBuddies) error
	Gossip(peer.Peer) (map[peer.Peer]peer.PeerInfo, error)
}
