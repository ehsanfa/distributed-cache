package gossip

import (
	"fmt"
)

const maxBuddyNum int = 1

type Buddy Peer

func (n *Node) setBuddy(b Buddy) {
	fmt.Printf("%s become buddy with %s \n", n.Name, b.Name)
	n.buddy = b
	buddyFound <- true
}

var buddyFound chan bool

func (n *Node) hasBuddy() bool {
	if n.buddy == (Buddy{}) {
		return false
	}
	return true
}

func (n *Node) checkForBuddies(noBuddyPeers []Peer) {
	if !thisNode.hasBuddy() && len(noBuddyPeers) > 0 {
		for _, peer := range noBuddyPeers {
			if peer == (Peer{thisNode.Name, thisNode.Port}) {
				continue
			}
			n.becomeBuddies(peer)
		}
	}
}

func (n *Node) BuddyRequest(peer Peer, resp *BuddyRequestResp) error {
	*resp = BuddyRequestResp{}
	if len(thisNode.buddyWith) <= maxBuddyNum {
		thisNode.buddyWith = append(thisNode.buddyWith, peer)
		resp.Res = true
	} else {
		resp.Res = false
	}
	return nil
}

func (n *Node) becomeBuddies(peer Peer) {
	c, err := dial(peer)
	if err != nil {
		panic(err)
	}
	var resp BuddyRequestResp
	p := Peer{thisNode.Name, thisNode.Port}
	c.Call("Node.BuddyRequest", p, &resp)
	if resp.Res {
		thisNode.setBuddy(Buddy(peer))
		n.touch()
		fmt.Println("touching node")
	}
}