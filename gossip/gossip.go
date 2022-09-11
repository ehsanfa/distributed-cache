package gossip

import (
	"fmt"
	"log"
	"time"
	"net"
	"net/rpc"
)

const gossipInterval = 5 * time.Second

var thisNode *Node

type GossipRequest struct {
	Info      map[Peer]PeerInfo
	Version   Version
	BuddyLook map[Peer]PeerInfo
}

func (g GossipRequest) GetInfo() map[Peer]PeerInfo {
	return g.Info
}

func (g GossipRequest) GetVersion() Version {
	return g.Version
}

func (g GossipRequest) GetBuddyLook() map[Peer]PeerInfo {
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
	var peer Peer
	go func() {
		for {
			if !n.hasBuddy() {
				peer = Peer(n.seeder)
			} else {
				peer = Peer(n.buddy)
			}
			select {
			case <-timer.C:
				n.doGossip(peer)
			case <-done:
				break
			}
		}
	}()
}

func (n *Node) doGossip(p Peer) {
	c, err := dial(p)
	if err != nil {
		panic(err)
	}
	var resp Response
	fmt.Println("buddy report", thisNode.Name, thisNode.buddy.Name)
	// TODO: Find a better way
	// var m []Peer
	// for peer, _ := range info {
	// 	m = append(m, peer)
	// }
	gossipRequest := GossipRequest{info, n.version, thisNode.noBuddyPeers}
	// fmt.Println("sending gossip request to ", n.buddy.Name, info)
	c.Call("Node.Gossip", gossipRequest, &resp)

	updateInfo(resp)

	thisNode.checkForBuddies()

	fmt.Println("info", info, thisNode.buddy.Name)
}

func updateInfo(g GossipMaterial) {
	/**
	 * This is a mess. Come back later and fix it FGS
	 */
	for peer, peerInfo := range g.GetInfo() {

		if _, ok := getInfo(peer); !ok {
			peer.track(peerInfo)
			continue
		}

		pi, _ := getInfo(peer)
		if pi.Version.compare(peerInfo.Version) >= 0 {
			// The message has nothing to give to us. Moving on
			continue
		}

		fmt.Println("SURPRISE. UPDATING", peer, pi, peerInfo)

		peer.track(peerInfo)

	}
	updateNoBuddyPeers(g)
}

func updateNoBuddyPeers(g GossipMaterial) {
	for peer, peerInfo := range g.GetInfo() {
		switch peerInfo.IsSomeonesBuddy {
		case true:
			if _, ok := thisNode.noBuddyPeers[peer]; ok {
				delete(thisNode.noBuddyPeers, peer)
			}
		case false:
			if _, ok := thisNode.noBuddyPeers[peer]; !ok {
				thisNode.noBuddyPeers[peer] = peerInfo
			}
		}
	}
}

func (n *Node) listen(done chan<- *Node) {
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
	done <- n
	
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
	updateInfo(req)
	// if !thisNode.hasBuddy() {
	// 	// thisNode.touch()
	// 	thisNode.noBuddyPeers[thisNode.getPeer()] = thisNode.getPeerInfo()
	// }
	thisNode.checkForBuddies()
	*resp = Response{
		Info: info, 
	}
	// fmt.Println("gossip result", resp)
	return nil
}