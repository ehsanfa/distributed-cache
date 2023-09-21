package node

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/gossip/gossip"
	"dbcache/cluster/info"
	"dbcache/cluster/network"
	"dbcache/cluster/network/rpc"
	"dbcache/cluster/peer"
	"log"
)

type Seeder struct {
	p peer.Peer
	n network.Network
}

func CreateSeederNode(
	p peer.Peer,
	cache cacher.Cache,
	i info.ClusterInfo,
	buff buffer.Buffer,
) (*Seeder, error) {
	n, err := rpc.CreateRpcNetwork(p, i, cache, buff)
	if err != nil {
		return nil, err
	}
	g, err := gossip.CreateGossipNetwork(n, i, nil, true, p)
	if err != nil {
		return nil, err
	}
	g.Start()
	return &Seeder{p, n}, nil
}

func (seeder *Seeder) Run() {
	log.Printf("Running seeder on %s port %d", seeder.p.Name(), seeder.n.Peer().Port())
	select {}
}
