package peer

import (
	"dbcache/cluster/partition"
	"encoding"
)

type Peer interface {
	Name() string
	WithPort
	Partition() partition.Partition
	IsSameAs(peer Peer) bool
	SetPartition(partition.Partition)
	SetPort(uint16)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type WithPort interface {
	Port() uint16
}

// func (n *Node) getPeersToShareBuffer() map[Peer]bool {
// 	peers := n.getPeersForPartitioning()
// 	if len(peers) == 0 {
// 		peers = n.convertInfoToPartitionPeers()
// 	}
// 	return peers
// }

// func (p *Peer) hasPartition(pi PeerInfo) bool {
// 	if pi.Partition == (partition.Partition{}) {
// 		return false
// 	}
// 	return true
// }

// func (p *Peer) track(i PeerInfo) {
// 	setInfo(*p, i)
// }
