package gossip

import (
	"fmt"
)

func (n *Node) introduce(ch chan<- bool) {
	fmt.Println("introduction")
	p := Peer(n.seeder)
	c, err := dial(p)
	if err != nil {
		panic(err)
	}
	var resp Response
	peer := Peer{n.Name, n.Port}
	c.Call("Node.Introduce", peer, &resp)
	ch <- true
	updateInfo(resp)
	c.Call("Node.SwitchBuddies", peer, &resp)
	thisNode.checkForBuddies()
}

func (n *Node) Introduce(peer Peer, resp *Response) error {
	// node.NewVersion()
	// if !thisNode.hasBuddy() {
	// 	b := Buddy(node)
	// 	thisNode.setBuddy(b)
	// 	return nil
	// }
	i := NewPeerInfo()
	thisNode.noBuddyPeers[peer] = i
	peer.track(i)
	*resp = Response{Info: info}
	fmt.Println(resp)
	fmt.Println(info, peer)
	updateInfo(resp)
	return nil
}