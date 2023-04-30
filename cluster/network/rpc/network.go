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
	"strings"
	"time"
)

type RpcNetwork struct {
	server *rpcServer
	info   info.ClusterInfo
	cache  cacher.Cache
	buffer buffer.Buffer
	p      peer.Peer
}

var hostNetwork *RpcNetwork

func CreateRpcNetwork(
	p peer.Peer,
	info info.ClusterInfo,
	cache cacher.Cache,
	buff buffer.Buffer,
) (*RpcNetwork, error) {
	n := &RpcNetwork{info: info, p: p, cache: cache, buffer: buff}
	hostNetwork = n
	c, err := n.serve()
	if err != nil {
		return nil, err
	}
	n.p = p.SetPort(c.Port())
	rpc.RegisterName(n.server.Node.rpcName(), n.server.Node)
	return n, nil
}

type RpcNode struct {
	client  *rpc.Client
	network *RpcNetwork
	p       peer.Peer
}

func (n *RpcNetwork) Connect(p peer.Peer, timeout time.Duration) (network.Node, error) {
	client, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()), timeout)
	if err != nil {
		return nil, err
	}
	return &RpcNode{client: rpc.NewClient(client), network: n, p: p}, nil
}

func (n *RpcNetwork) Peer() peer.Peer {
	return n.p
}

func (n *RpcNetwork) Kill() {
	n.server.kill <- true
}

func (n *RpcNode) rpcName() string {
	name := strings.Replace(n.p.Name(), ".", "", -1)
	return fmt.Sprintf("%s%d", name, n.p.Port())
}

func (n *RpcNode) rpcAction(action string) string {
	return fmt.Sprintf("%s.%s", n.rpcName(), action)
}
