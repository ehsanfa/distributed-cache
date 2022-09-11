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
	thisNode.checkForBuddies(resp.BuddyLook)
}

func (n *Node) Introduce(node Peer, resp *Response) error {
	thisNode.noBuddyPeers[node] = true
	node.NewVersion()
	updateInfo(node)
	// if !thisNode.hasBuddy() {
	// 	b := Buddy(node)
	// 	thisNode.setBuddy(b)
	// 	return nil
	// }
	*resp = Response{Info: info, BuddyLook: thisNode.noBuddySlice()}
	return nil
}