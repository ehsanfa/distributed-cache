package cluster

import (
	"fmt"
	"testing"
	// "time"
)

type node struct {
	endSignal chan bool
	n         Node
}

func TestCluster(t *testing.T) {
	s := node{make(chan bool), CreateNode(true)}
	s.n.Initialize(s.endSignal)
	seeder := CreateSeeder(s.n.Peer.Name, fmt.Sprint(s.n.Peer.Port))
	n1 := node{make(chan bool), CreateNode(false)}
	n1.n.SetSeeder(seeder)
	n1.n.Initialize(n1.endSignal)
	// time.Sleep(2 * time.Second)
	t.Error(info)
}