package info

import (
	"dbcache/cluster/peer"
)

type ClusterInfo interface {
	Get(peer.Peer) (peer.PeerInfo, bool)
	Add(peer.Peer, peer.PeerInfo)
	Upsert(peer.Peer, peer.PeerInfo)
	All() map[peer.Peer]peer.PeerInfo
	AllAlive() map[peer.Peer]peer.PeerInfo
	IsPeerKnown(peer.Peer) bool
	IsPeerAlive(peer.Peer) bool
	List() []peer.Peer
	Remove(peer.Peer)
	Replace(map[peer.Peer]peer.PeerInfo)
	Update(map[peer.Peer]peer.PeerInfo)
	MarkAsDead(peer.Peer) error
}

// func getPartitionPeers(p partition.Partition) []Peer {
// 	var peers []Peer
// 	for peer, pi := range info {
// 		if p == pi.Partition && peer != *thisNode.getPeer() {
// 			peers = append(peers, peer)
// 		}
// 	}
// 	return peers
// }

// func setInfo(peer Peer, i PeerInfo) {
// 	infoMutex.Lock()
// 	info[peer] = i
// 	infoMutex.Unlock()
// }

// type ShareInfoRequest struct{}
// type ShareInfoResponse struct {
// 	Info       map[Peer]bool
// 	Partitions map[partition.Partition]map[Peer]bool
// }

// func getInfoToShare() map[Peer]bool {
// 	peers := make(map[Peer]bool)
// 	for peer, _ := range info {
// 		if peer.isAlive() {
// 			peers[peer] = true
// 		}
// 	}
// 	return peers
// }

// func (n *Node) ShareInfo(req ShareInfoRequest, resp *ShareInfoResponse) error {
// 	*resp = ShareInfoResponse{
// 		Info:       getInfoToShare(),
// 		Partitions: assignedPartitions,
// 	}
// 	return nil
// }
