package gossip

import (
	"strconv"
)

func CreateSeeder(name, port string) Seeder {
	s := Seeder{}
	s.Name = name
	p, _ := strconv.Atoi(port)
	s.Port = Port(p)
	return s
}

func CreateNode(isSeeder bool) Node {
	p := Peer{}
	buddies := make(map[Peer]bool)
	return Node{isSeeder: isSeeder, Peer: &p, buddies: buddies}
}