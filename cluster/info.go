package cluster

import (
	"sync"
	partition "dbcache/partitioning"
)

var info map[Peer]PeerInfo
var infoMutex sync.RWMutex

func getInfoList() []Peer {
	var peers []Peer
	for p, _ := range info {
		peers = append(peers, p)
	}
	return peers
}

func getPartitionPeers(p partition.Partition) []Peer {
	var peers []Peer
	for peer, pi := range info {
		if p == pi.Partition && peer != *thisNode.getPeer() {
			peers = append(peers, peer)
		}
	}
	return peers
}

func getInfo(peer Peer) (PeerInfo, bool){
	infoMutex.RLock()
	v, ok := info[peer]
	infoMutex.RUnlock()
	return v, ok
}

func setInfo(peer Peer, i PeerInfo) {
	infoMutex.Lock()
	info[peer] = i
	infoMutex.Unlock()
}

type ShareInfoRequest struct {}
type ShareInfoResponse struct {
	Info       map[Peer]bool
	Partitions map[partition.Partition]map[Peer]bool
}

func getInfoToShare() map[Peer]bool {
	peers := make(map[Peer]bool)
	for peer, _ := range info {
		if peer.isAlive() {
			peers[peer] = true
		}
	}
	return peers
}

func (n *Node) ShareInfo(req ShareInfoRequest, resp *ShareInfoResponse) error {
	*resp = ShareInfoResponse{
		Info: getInfoToShare(),
		Partitions: assignedPartitions,
	}
	return nil
}