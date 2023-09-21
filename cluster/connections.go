package cluster

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

func (n *Node) dial(p Peer) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
}

func (n *Node) getConnection(p Peer) (*rpc.Client, error) {
	// return n.dial(p)
	var mu sync.RWMutex
	var conn *rpc.Client
	var err error
	mu.RLock()
	conn, ok := n.connections[p]
	mu.RUnlock()
	if !ok {
		fmt.Println("no cache, creating connection")
		conn, err = n.dial(p)
		if err != nil {
			return nil, err
		}
		mu.Lock()
		n.connections[p] = conn
		mu.Unlock()
	}
	return conn, nil
}

func (n *Node) listen(done chan<- *Node) {
	node := new(Node)
	rpc.Register(node)

	var listener net.Listener
	var err error

	if n.isSeeder {
		listener, err = net.Listen("tcp", "0.0.0.0:7000")
	} else {
		listener, err = net.Listen("tcp", "0.0.0.0:")
	}

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	n.setPort(listener)
	done <- n

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
