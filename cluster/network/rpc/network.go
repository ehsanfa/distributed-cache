package rpc

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/network"
	"dbcache/cluster/peer"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type RpcNetwork struct {
	server *rpcServer
	info   info.ClusterInfoProvider
	cache  cacher.Cache
	buffer buffer.Buffer
	p      peer.Peer
}

var hostNetwork *RpcNetwork

func CreateRpcNetwork(
	p peer.Peer,
	info info.ClusterInfoProvider,
	cache cacher.Cache,
	buff buffer.Buffer,
) (*RpcNetwork, error) {
	n := &RpcNetwork{info: info, p: p, cache: cache, buffer: buff}
	hostNetwork = n
	rpc.Register(n)
	c, err := n.serve()
	if err != nil {
		return nil, err
	}
	n.p = p.SetPort(c.Port())
	return n, nil
}

type RpcNode struct {
	client  *rpc.Client
	network *RpcNetwork
}

func (n *RpcNetwork) Connect(p peer.Peer, timeout time.Duration) (network.Node, error) {
	client, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()), timeout)
	if err != nil {
		return nil, err
	}
	node := &RpcNode{client: rpc.NewClient(client), network: n}
	rpc.Register(node)
	return node, nil
}

func (n *RpcNetwork) Peer() peer.Peer {
	return n.p
}

func (n *RpcNetwork) Kill() {
	n.server.kill <- true
}
