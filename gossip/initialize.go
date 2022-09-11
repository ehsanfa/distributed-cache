package gossip

import (
	"os"
	"fmt"
)

func (n *Node) Initialize() {
	info = make(map[Peer]PeerInfo)
	buddyFound = make(chan bool)
	if n.isSeeder {
		n.newVersion()
	}
	n.noBuddyPeers = make(map[Peer]PeerInfo)
	thisNode = n
	nodename, err := os.Hostname()
	if err != nil {
		panic("Unable to set host name")
	}
	n.Name = nodename
	// if len(info) > 0 {
	// 	buddyNode := info[len(info)-1]
	// 	n.buddy = Buddy{buddyNode.Name, buddyNode.Port}
	// }
	done := make(chan *Node)
	go n.listen(done)
	n = <-done
	close(done)

	en := Peer{n.Name, n.Port}
	peerInfo := NewPeerInfo()
	if n.isSeeder {
		thisNode.noBuddyPeers[en] = peerInfo
	}
	setInfo(en, peerInfo)

	if !n.isSeeder {
		introduced := make(chan bool)
		go n.introduce(introduced)
		<-introduced
		close(introduced)
	}
	
	<-buddyFound
	fmt.Println("buddy found", thisNode.Name, thisNode.buddy)
	n.gossip()
}