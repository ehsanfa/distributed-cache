package gossip

import (
	innercircle "dbcache/cluster/gossip/inner_circle"
	"dbcache/cluster/info"
	"dbcache/cluster/network"
	"dbcache/cluster/peer"
	"time"
)

type Seeder peer.Peer

const maxBuddyNum = 3
const interval = 1 * time.Second
const timeout = 5 * time.Second

type Gossip struct {
	network      network.Network
	info         info.ClusterInfo
	inner_cirlce innercircle.InnerCircle
}

func CreateGossipNetwork(network network.Network, info info.ClusterInfo, seeder peer.Peer) *Gossip {
	inner_cirle := innercircle.CreateInMemoryBuddies(maxBuddyNum)
	inner_cirle.Add(seeder)
	return &Gossip{network, info, inner_cirle}
}

func (g *Gossip) Start() {
	timer := time.NewTicker(interval)
	go func() {
		for {
			<-timer.C
			g.spawn()
		}
		// for {
		// select {
		// case <-timer.C:
		// 	g.spawn()
		// }
		// }
	}()
}

func (g *Gossip) spawn() {
	for peer := range g.inner_cirlce.All() {
		go g.gossip(peer)
	}
}

func (g *Gossip) gossip(p peer.Peer) {
	node, err := g.network.Connect(p, timeout)
	if err != nil {
		g.gossipFailed(p, err)
	}
	info, err := node.GetClusterInfo()
	if err != nil {
		g.gossipFailed(p, err)
	} else {
		g.info.Update(info)
	}
}

func (g *Gossip) gossipFailed(p peer.Peer, err error) {
	panic(err)
}
