package cluster

import (
	"fmt"
)

const maxBuddyNum int = 2

func (n *Node) addBuddy(b Peer) {
	n.mu.Lock()
	n.buddies[b] = true
	n.mu.Unlock()
}

func (n *Node) getBuddies() map[Peer]bool {
	return n.buddies
}

func (n *Node) buddyCount() int {
	return len(n.buddies)
}

func (n *Node) hasBuddy() bool {
	return n.buddyCount() > 0
}

func (n *Node) isBuddyWith(peer Peer) bool {
	for p, _ := range n.buddies {
		if p == peer {
			return true
		}
	}
	return false
}

func (n *Node) canAcceptBuddyRequest(peer Peer) bool {
	if n.isBuddyWith(peer) {
		return false
	}
	return true
}

func (n *Node) checkForBuddies() {
	peers := getInfoList()
	if n.buddyCount() < maxBuddyNum {
		for _, peer := range peers {
			thisPeer := n.getPeer()
			if peer.isSame(*thisPeer) || !peer.isAlive() || n.isBuddyWith(peer) {
				continue
			}
			fmt.Println("buddy request", peer)
			if n.becomeBuddies(peer) {
				break
			} 
		}
	}
}

func (n *Node) BuddyRequest(peer Peer, resp *BuddyRequestResp) error {
	*resp = BuddyRequestResp{true}

	if !thisNode.canAcceptBuddyRequest(peer) {
		resp.Res = false
	}

	return nil
}

func (n *Node) becomeBuddies(peer Peer) bool {
	c, err := n.dial(peer)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	var resp BuddyRequestResp
	p := n.getPeer()
	c.Call("Node.BuddyRequest", p, &resp)
	if resp.Res {
		thisNode.addBuddy(peer)
		return true
	}
	return false
}