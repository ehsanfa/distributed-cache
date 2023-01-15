package info

import (
	"dbcache/cluster/peer"
)

type ClusterInfo interface {
	Add(p peer.Peer, i peer.PeerInfo)
	All() map[peer.Peer]peer.PeerInfo
	IsPeerKnown(p peer.Peer) bool
	IsPeerAlive(p peer.Peer) bool
	List() []peer.Peer
	Remove(p peer.Peer)
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
