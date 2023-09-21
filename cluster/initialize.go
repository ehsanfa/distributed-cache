package cluster

import (
	"os"
	// "fmt"
	// "time"
	"dbcache/cluster/buffer"
	"net/rpc"
	// ll "github.com/ehsanfa/linked-list"
)

func (n *Node) Initialize(endSignal chan bool) {
	// n.cacheVersions = make(map[string]cacher.CacheVersion)
	n.connections = make(map[Peer]*rpc.Client)
	n.buffer = buffer.CreateInMemoryBuffer()

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
		n.Introduce()
	}

	// time.Sleep(time.Second * 5)
	go n.startGossiping(endSignal)
	go n.startCleaningBuffer()
	// go n.reportCount()
	n.syncCacheWithPeers()
}
