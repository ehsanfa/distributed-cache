package node

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/gossip/gossip"
	"dbcache/cluster/info"
	"dbcache/cluster/network/rpc"
	"dbcache/cluster/peer"
	"log"
	"time"
)

type Cacher struct {
	p      peer.Peer
	seeder peer.Peer
}

func CreateCacherNode(
	p peer.Peer,
	seeder peer.Peer,
	cache cacher.Cache,
	i info.ClusterInfo,
	buff buffer.Buffer,
) (*Cacher, error) {
	if seeder == nil {
		panic("Seeder info required")
	}
	n, err := rpc.CreateRpcNetwork(p, i, cache, buff)
	if err != nil {
		return nil, err
	}
	g, err := gossip.CreateGossipNetwork(n, i, seeder, false, p)
	if err != nil {
		return nil, err
	}
	node, err := n.Connect(seeder, 10*time.Second)
	if err != nil {
		return nil, err
	}
	err = node.Introduce(peer.Cacher, n.Peer())
	if err != nil {
		return nil, err
	}
	g.Start()
	return &Cacher{p, seeder}, nil
}

func (cacher *Cacher) Run() {
	log.Printf("Running cacher on %s port %d", cacher.p.Name(), cacher.p.Port())
	select {}
}
