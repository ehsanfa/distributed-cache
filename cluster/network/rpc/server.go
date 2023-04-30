package rpc

import (
	"dbcache/cluster/peer"
	"fmt"
	"net"
	"net/rpc"
)

type rpcServer struct {
	Node *RpcNode
	port uint16
	kill chan bool
}

func (n *RpcNetwork) serve() (peer.WithPort, error) {
	ch := make(chan connRes)
	kill := make(chan bool)
	node := &RpcNode{p: n.p}
	n.server = &rpcServer{Node: node, kill: kill}

	var err error
	go n.server.initialize(n.p, ch)
	res := <-ch
	if res.err != nil {
		return nil, res.err
	}

	n.server.port = res.port
	return n.server, err
}

type connRes struct {
	success bool
	port    uint16
	err     error
}

func (s *rpcServer) initialize(p peer.Peer, channel chan<- connRes) {
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
	port := extractPort(listener)
	channel <- connRes{true, port, err}
	s.Node.p = s.Node.p.SetPort(port)

	server := rpc.NewServer()
	server.RegisterName(s.Node.rpcName(), s.Node)

	for {
		conn, err := listener.Accept()

		if err != nil {
			panic(err)
		}
		go server.ServeConn(conn)
		go func() {
			<-s.kill
			conn.Close()
			listener.Close()
		}()
	}
}

func extractPort(l net.Listener) uint16 {
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func (c *rpcServer) Port() uint16 {
	return c.port
}
