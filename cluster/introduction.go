package cluster

import (
	"fmt"
)

func (n *Node) introduce() {
	fmt.Println("introduction")
	p := n.getSeeder().getPeer()
	c, err := n.dial(p)
	if err != nil {
		panic(err)
	}
	var resp Response
	c.Call("Node.Introduce", n.getPeer(), &resp)
	updateInfo(resp)
	n.partition = resp.Partition
}

func (n *Node) Introduce(peer Peer, resp *Response) error {
	i := NewPeerInfo()
	thisNode.checkForBuddies()
	i.assignPartition(peer)
	i.touch()
	peer.track(i)
	fmt.Println("partitions", peer, i.Partition)
	*resp = Response{Info: info, Partition: i.Partition}
	return nil
}

func (n *Node) syncCacheWithPeers() {
	peers := getPartitionPeers(n.partition)
	for _, p := range peers {
		n.askForCache(p)
	}
}