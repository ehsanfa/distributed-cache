package gossip

import (
	"fmt"
	"log"
	"time"
	"net"
	"net/rpc"
)

const gossipInterval = 5 * time.Second

var info map[Peer]PeerInfo

var thisNode *Node

type GossipRequest struct {
	Info      []Peer
	Version   Version
	BuddyLook []Peer
}

func (g GossipRequest) GetInfo() []Peer {
	return g.Info
}

func (g GossipRequest) GetVersion() Version {
	return g.Version
}

func (g GossipRequest) GetBuddyLook() []Peer {
	return g.BuddyLook
}

/*
 - When a node gets initialized, it makes a request to its seeder to get its info
 - Seeder sends its list of known nodes and the initialized node picks one to be buddy with
 - If the seeder doesn't have a buddy yet itself, it becomes the node's buddy
 - Otherwise it spreads the word that this nodes needs a buddy
 - If the messages get to a node that could get the newjoiner as its buddy, it makes a
   request to the newjoiner's BecomeBuddy function and becomes it buddy
 - There should be max limit on how many buddies a node can have
*/

func (n *Node) gossip() {
	timer := time.NewTicker(gossipInterval)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-timer.C:
				n.doGossip()
			case <-done:
				break
			}
		}
	}()
}

func (n *Node) doGossip() {
	p := Peer(n.buddy)
	c, err := dial(p)
	if err != nil {
		panic(err)
	}
	var resp Response
	// TODO: Find a better way
	var m []Peer
	for peer, _ := range info {
		m = append(m, peer)
	}
	gossipRequest := GossipRequest{m, n.version, thisNode.noBuddySlice()}
	// fmt.Println("sending gossip request to ", n.buddy.Name, info)
	c.Call("Node.Gossip", gossipRequest, &resp)

	fmt.Println("buddy lookup", resp.BuddyLook)

	n.update(resp)

	thisNode.checkForBuddies(thisNode.noBuddySlice())

	fmt.Println("info", info, thisNode.buddy.Name)
}

func updateInfo(peer Peer) {

}

func (n *Node) update(g GossipMaterial) {
	/**
	 * TODO: Find a better approach
	 * 
	*/
	// if n.version.compare(g.GetVersion()) >= 0 {
		for _, en := range g.GetInfo() {
			info[en] = true
		}
		for _, peer := range g.GetBuddyLook() {
			if peer == n.getPeer() && n.hasBuddy() {
				delete(thisNode.noBuddyPeers, peer)
				continue
			}
			if _, ok := thisNode.noBuddyPeers[peer]; !ok {
				thisNode.noBuddyPeers[peer] = true
			}
		}
		n.version.replace(g.GetVersion())
		// fmt.Println("updating from resp", g, info)
	// }
}

func (n *Node) listen(done chan<- bool) {
	node := new(Node)
    rpc.Register(node)

    var listener net.Listener
    var err error

    if n.isSeeder {
    	listener, err = net.Listen("tcp", "0.0.0.0:7000")
    } else {
    	listener, err = net.Listen("tcp", "0.0.0.0:")
    }
	
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	
	n.setPort(listener)
	done <- true
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (n *Node) Gossip(req GossipRequest, resp *Response) error {
	/**
	 * TODO: change the way buddy is chosen
	 * */
	// if !thisNode.hasBuddy() {
	// 	buddyNode := Buddy{node.Name, node.Port}
	// 	thisNode.buddy = buddyNode
	// }
	// if thisNode.isSeeder {
		fmt.Println("nodes with no buddy", thisNode.noBuddyPeers)
	// }
	n.update(req)
	if !thisNode.hasBuddy() {
		thisNode.noBuddyPeers[thisNode.getPeer()] = true
	}
	thisNode.checkForBuddies(thisNode.noBuddySlice())
	*resp = Response{
		Info: info, 
		BuddyLook: thisNode.noBuddySlice(), 
		Version: thisNode.version,
	}
	// fmt.Println("gossip result", resp)
	return nil
}