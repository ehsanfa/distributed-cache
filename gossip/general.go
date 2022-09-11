package gossip

import (
	"fmt"
	"net/rpc"
	"strconv"
)

func dial(p Peer) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
}

func CreateSeeder(name, port string) Seeder {
	s := Seeder{}
	s.Name = name
	p, _ := strconv.Atoi(port)
	s.Port = Port(p)
	return s
}

func CreateNode(isSeeder bool) Node {
	return Node{isSeeder: isSeeder}
}