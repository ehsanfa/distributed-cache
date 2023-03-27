package rpc

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/info"
	"dbcache/cluster/network"
	"dbcache/cluster/peer"
	"fmt"
	"net/rpc"
)

type RpcNetwork struct {
	server *rpcServer
	info   info.ClusterInfoProvider
	cache  network.CacheProvider
	buffer buffer.Buffer
	p      peer.Peer
}

var hostNetwork *RpcNetwork

func CreateRpcNetwork(
	p peer.Peer,
	info info.ClusterInfoProvider,
	cache network.CacheProvider,
	buff buffer.Buffer,
) (*RpcNetwork, error) {
	n := &RpcNetwork{info: info, p: p, cache: cache, buffer: buff}
	hostNetwork = n
	rpc.Register(n)
	c, err := n.Serve(p)
	if err != nil {
		return nil, err
	}
	p.SetPort(c.Port())
	return n, nil
}

type RpcNode struct {
	client  *rpc.Client
	network *RpcNetwork
}

func (n *RpcNetwork) Connect(p peer.Peer) (network.Node, error) {
	s, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()))
	if err != nil {
		return nil, err
	}
	node := &RpcNode{client: s, network: n}
	rpc.Register(node)
	return node, nil
}

func (n *RpcNetwork) Kill() {
	n.server.kill <- true
}
