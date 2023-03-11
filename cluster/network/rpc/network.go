package network

import (
	"dbcache/cluster/network"
	"dbcache/cluster/peer"
	"fmt"
	"net/rpc"
)

type RpcNetwork struct {
	server *rpcServer
}

func CreateRpcNetwork() *RpcNetwork {
	n := &RpcNetwork{}
	rpc.Register(n)
	return n
}

type RpcNode struct {
	client *rpc.Client
}

func (n *RpcNetwork) Connect(p peer.Peer) (network.Node, error) {
	s, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()))
	if err != nil {
		return nil, err
	}
	node := &RpcNode{client: s}
	return node, nil
}

func (n *RpcNetwork) Kill() {
	n.server.kill <- true
}
