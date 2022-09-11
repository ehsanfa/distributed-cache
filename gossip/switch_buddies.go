package gossip

import (
	"fmt"
)

func (n *Node) isStuck() {
	p := Peer(n.seeder)
	c, err := dial(p)
	if err != nil {
		panic(err)
	}
	var resp Response
	peer := Peer{n.Name, n.Port}
	c.Call("Node.SwitchBuddies", peer, &resp)
	updateInfo(resp)
	thisNode.checkForBuddies()
}

func (n *Node) switchBuddies(target Peer) {
	oldBuddy := Peer(thisNode.buddy)
	if n.hasBuddy() {
		// seeder switches buddy with the new joiner
		i, _ := getInfo(oldBuddy)
		i.touch()
		setInfo(oldBuddy, i)
		i.IsSomeonesBuddy = false
	}
	n.buddy = Buddy(target)
	fmt.Println("seeder switched buddies", oldBuddy, target)
}

func (n *Node) SwitchBuddies(stuckNode Peer, resp *Response) error {
	thisNode.switchBuddies(stuckNode)
	*resp = Response{Info: info}
	return nil
}