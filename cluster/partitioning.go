package cluster

import (
	partition "dbcache/partitioning"
	"fmt"
	"sync"
)

var partitions []partition.Partition
var assignedPartitions map[partition.Partition]map[Peer]bool

func (n *Node) getPeersForPartitioning() map[Peer]bool {
	p := n.partition
	peers, ok := assignedPartitions[p]
	fmt.Println("assigned partitions", p, assignedPartitions)
	if !ok {
		return map[Peer]bool{}
	}
	return peers
}

func (n *Node) convertInfoToPartitionPeers() map[Peer]bool {
	peers := make(map[Peer]bool)
	for _, p := range getInfoList() {
		peers[p] = true
	}
	return peers
}

func (peerInfo *PeerInfo) assignPartition(peer Peer) {
	if len(partitions) == 0 {
		partitions = partition.Initialize(float64(minPartitions))
	}

	if assignedPartitions == nil {
		assignedPartitions = make(map[partition.Partition]map[Peer]bool)
	}

	smallestPartition := partitions[0]
	for _, p := range partitions {
		if _, ok := assignedPartitions[p]; !ok {
			smallestPartition = p
			break
		}
		if len(assignedPartitions[p]) <= len(assignedPartitions[smallestPartition]) {
			smallestPartition = p
		}
	}

	peerInfo.Partition = smallestPartition
	addToParitionsInfo(peer, *peerInfo)
}

func addToParitionsInfo(peer Peer, peerInfo PeerInfo) {
	var mu sync.Mutex
	mu.Lock()
	if assignedPartitions == nil {
		assignedPartitions = make(map[partition.Partition]map[Peer]bool)
	}
	if assignedPartitions[peerInfo.Partition] == nil {
		assignedPartitions[peerInfo.Partition] = make(map[Peer]bool)
	}
	assignedPartitions[peerInfo.Partition][peer] = true
	defer mu.Unlock()
}
