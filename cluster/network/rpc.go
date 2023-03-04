package network

import (
	"dbcache/cluster/info"
	"dbcache/cluster/peer"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type RpcServer struct {
	Peer peer.Peer
	port uint16
	info info.ClusterInfoProvider
}

func (c *RpcServer) connect(p peer.Peer) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()))
}

func (c *RpcServer) GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error) {
	client, err := c.connect(c.Peer)
	if err != nil {
		return nil, err
	}
	resp := new(ClusterInfoResponse)
	client.Call("RpcServer.RpcGetClusterInfo", struct{}{}, &resp)
	return resp.ClusterInfo, nil
}

type ClusterInfoResponse struct {
	ClusterInfo map[peer.Peer]peer.PeerInfo
}

func (c *RpcServer) RpcGetClusterInfo(p struct{}, resp *ClusterInfoResponse) error {
	*resp = ClusterInfoResponse{c.info.GetClusterInfo()}
	return nil
}

func CreateRpcServer(p peer.Peer, info info.ClusterInfoProvider) (*RpcServer, error) {
	ch := make(chan connRes)
	var err error
	go initialize(p, ch)
	res := <-ch
	if res.err != nil {
		return nil, res.err
	}

	c := &RpcServer{Peer: p, port: res.port, info: info}
	rpc.Register(c)
	return c, err
}

type connRes struct {
	success bool
	port    uint16
	err     error
}

func initialize(p peer.Peer, channel chan<- connRes) {
	var listener net.Listener
	var err error

	if p.Port() == 0 {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:", p.Name()))
	} else {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", p.Name(), p.Port()))
	}

	if err != nil {
		channel <- connRes{false, 0, err}
		return
	}
	defer listener.Close()
	channel <- connRes{true, extractPort(listener), err}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (c *RpcServer) Ping() (bool, error) {
	var resp PingResponse
	rc, err := c.connect(c.Peer)
	if err != nil {
		return false, err
	}
	err = rc.Call("RpcServer.RpcPing", PingRequest{}, &resp)
	if err != nil {
		return false, err
	}
	return true, nil
}

type PingRequest struct{}

type PingResponse struct{}

func (c *RpcServer) RpcPing(req PingRequest, resp *PingResponse) error {
	return nil
}

func extractPort(l net.Listener) uint16 {
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func (c *RpcServer) Port() uint16 {
	return c.port
}
