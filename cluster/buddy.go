package cluster

import (
	"fmt"
)

const maxBuddyNum int = 2
var alreadyRequested map[Peer]bool

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
			_, ar := alreadyRequested[peer]
			thisPeer := n.getPeer()
			if peer.isSame(*thisPeer) || !peer.isAlive() || n.isBuddyWith(peer) || ar {
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

func (n *Node) askForCache(peer Peer) {
	c, err := n.dial(peer)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	var resp ShareCacheResposne
	req := ShareCacheRequest{}
	c.Call("Node.ShareCache", req, &resp)
	if len(resp.Cache) > 0 {
		thisNode.cache = resp.Cache
	}
}

func (n *Node) becomeBuddies(peer Peer) bool {
	c, err := n.dial(peer)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer c.Close()
	var resp BuddyRequestResp
	p := n.getPeer()
	c.Call("Node.BuddyRequest", p, &resp)
	if resp.Res {
		thisNode.addBuddy(peer)
		thisNode.askForCache(peer)
		fmt.Println("became buddy with", n.getPeer(), peer)
		return true
	}
	if alreadyRequested == nil {
		alreadyRequested = make(map[Peer]bool)
	}
	alreadyRequested[peer] = true
	return false
}