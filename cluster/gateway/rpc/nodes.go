package rpc

import (
	"dbcache/cluster/network"
	"time"
)

func (g *RpcGateway) getNextCacher() (network.Node, error) {
	for p, pi := range g.clusterInfo.AllAlive() {
		if pi.IsCacher() {
			n, e := g.network.Connect(p, 1*time.Second)
			if e != nil {
				return nil, e
			}
			return n, nil
		}
	}
	return nil, nil
}
