package node

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	rpcGateway "dbcache/cluster/gateway/rpc"
	"dbcache/cluster/gossip/gossip"
	"dbcache/cluster/info"
	"dbcache/cluster/network/rpc"
	"dbcache/cluster/peer"
	"log"
	"time"
)

type Gateway struct {
	p      peer.Peer
	seeder peer.Peer
}

func CreateGatewayNode(
	p peer.Peer,
	seeder peer.Peer,
	cache cacher.Cache,
	i info.ClusterInfo,
	buff buffer.Buffer,
) (*Gateway, error) {
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
	err = node.Introduce(peer.Gateway, n.Peer())
	if err != nil {
		return nil, err
	}
	gateway := rpcGateway.CreateRpcGateway(i, n)
	go gateway.Serve()
	g.Start()
	return &Gateway{p, seeder}, nil
}

func (gateway *Gateway) Run() {
	log.Printf("Running cacher on %s port %d", gateway.p.Name(), gateway.p.Port())
	select {}
}
