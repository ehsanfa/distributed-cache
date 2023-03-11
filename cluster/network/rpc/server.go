package network

import (
	"dbcache/cluster/info"
	"dbcache/cluster/peer"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type rpcServer struct {
	Peer peer.Peer
	port uint16
	info info.ClusterInfoProvider
	kill chan bool
}

func (n *RpcNetwork) Serve(p peer.Peer, info info.ClusterInfoProvider) (peer.WithPort, error) {
	ch := make(chan connRes)
	kill := make(chan bool)
	c := &rpcServer{Peer: p, info: info, kill: kill}
	n.server = c

	var err error
	go n.server.initialize(p, ch)
	res := <-ch
	if res.err != nil {
		return nil, res.err
	}

	c.port = res.port
	return c, err
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
	channel <- connRes{true, extractPort(listener), err}

	// for {
	conn, err := listener.Accept()

	if err != nil {
		log.Print(err)
		return
		// break
	}
	go rpc.ServeConn(conn)
	go func() {
		<-s.kill
		conn.Close()
		listener.Close()
	}()
	// }
}

func extractPort(l net.Listener) uint16 {
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func (c *rpcServer) Port() uint16 {
	return c.port
}
