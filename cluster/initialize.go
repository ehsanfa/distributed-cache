package cluster

import (
	"os"
	// "fmt"
	// "time"
	"net/rpc"
)

func (n *Node) Initialize() {
	info = make(map[Peer]PeerInfo)
	n.cache = make(map[string]string)
	n.connections = make(map[Peer]*rpc.Client)

	thisNode = n

	nodename, err := os.Hostname()
	if err != nil {
		panic("Unable to set host name")
	}
	n.setName(nodename)
	
	done := make(chan *Node)
	go n.listen(done)
	select {
	case n = <-done:
	}
	close(done)

	peer := n.getPeer()
	peerInfo := NewPeerInfo()

	peer.track(peerInfo)

	if !n.isSeeder {
		n.introduce()
	}

	// time.Sleep(time.Second * 5)
	endSignal := make(chan bool)
	n.startGossiping(endSignal)
	<-endSignal
}