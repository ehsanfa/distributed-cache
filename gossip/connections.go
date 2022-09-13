package gossip

import (
	"fmt"
	"net"
	"net/rpc"
	"log"
)

func (n *Node) dial(p Peer) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
	// if n.connections == nil {
	// 	n.connections = make(map[Peer]*rpc.Client)
	// }

	// if _, ok := n.connections[p]; !ok {
	// 	conn, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// n.connections[p] = conn
	// }

	// return n.connections[p], nil
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