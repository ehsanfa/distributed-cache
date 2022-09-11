package gossip

import (
	"fmt"
	"sync"
)

const maxBuddyNum int = 1

type Buddy Peer

func (n *Node) setBuddy(b Buddy) {
	fmt.Printf("%s became buddy with %s \n", n.Name, b.Name)
	n.buddy = b
	peer := Peer(b)
	i, _ := getInfo(peer)
	i.touch()
	i.IsSomeonesBuddy = true
	// n.noBuddyPeers[peer] = i
	peer.track(i)
	// fmt.Println("info for the buddy", peer, i, info)
	buddyFound <- true
}

var buddyFound chan bool

func (n *Node) hasBuddy() bool {
	if n.buddy == (Buddy{}) {
		return false
	}
	return true
}

var retry int8

func (n *Node) checkForBuddies() {
	// fmt.Println("checking for buddies", n.hasBuddy(), n.noBuddyPeers)
	if !n.hasBuddy(){
		for peer, _ := range n.noBuddyPeers {
			if peer == (Peer{thisNode.Name, thisNode.Port}) {
				continue
			}
			fmt.Println("buddy request", peer)
			if n.becomeBuddies(peer) {
				break
			} 
		}
	}
}

var mu sync.Mutex

func (n *Node) BuddyRequest(peer Peer, resp *BuddyRequestResp) error {
	*resp = BuddyRequestResp{}
	mu.Lock()
	defer mu.Unlock()
	isBuddyWith := false
	for _, p := range thisNode.buddyWith {
		if p == thisNode.getPeer() {
			isBuddyWith = true
		}
	}
	fmt.Println("is already buddy with", thisNode.Name, peer.Name, isBuddyWith)

	if len(thisNode.buddyWith) < maxBuddyNum && !isBuddyWith{
		thisNode.buddyWith = append(thisNode.buddyWith, peer)
		resp.Res = true
	} else {
		resp.Res = false
	}
	return nil
}

func (n *Node) becomeBuddies(peer Peer) bool {
	c, err := dial(peer)
	if err != nil {
		panic(err)
	}
	var resp BuddyRequestResp
	p := Peer{thisNode.Name, thisNode.Port}
	c.Call("Node.BuddyRequest", p, &resp)
	if resp.Res {
		thisNode.setBuddy(Buddy(peer))
		fmt.Println("touching node")
		return true
	}
	return false
}