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
}

func (n *Node) Introduce(peer Peer, resp *Response) error {
	i := NewPeerInfo()
	peer.track(i)
	thisNode.checkForBuddies()
	*resp = Response{Info: info}
	return nil
}