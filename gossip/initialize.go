package gossip

import (
	"os"
	"fmt"
)

func (n *Node) Initialize() {
	info = make(map[Peer]bool)
	buddyFound = make(chan bool)
	if n.isSeeder {
		n.newVersion()
	}
	n.noBuddyPeers = make(map[Peer]bool)
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
	done := make(chan bool)
	go n.listen(done)
	<-done
	close(done)

	if !n.isSeeder {
		introduced := make(chan bool)
		go n.introduce(introduced)
		<-introduced
		close(introduced)
	}

	en := Peer{n.Name, n.Port}
	if n.isSeeder {
		thisNode.noBuddyPeers[en] = true
	}
	info[en] = true
	<-buddyFound
	fmt.Println("buddy found", thisNode.Name, thisNode.buddy)
	n.gossip()
}