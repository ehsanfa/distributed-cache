package network

import (
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"testing"
)

func TestWithRandomPort(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", 0, &part)
	i := info.CreateInMemoryClusterInfo()
	network := CreateRpcNetwork()
	c, err := network.Serve(p, i)
	p.SetPort(c.Port())

	if err != nil {
		t.Error(err)
	}

	n, err := network.Connect(p)
	if err != nil {
		t.Error(err)
	}
	res, err := n.Ping()
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("node should be alive")
	}
	network.Kill()
}

func TestWithSpecificPort(t *testing.T) {
	port := uint16(63447)
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", port, &part)
	i := info.CreateInMemoryClusterInfo()
	network := CreateRpcNetwork()
	c, err := network.Serve(p, i)
	p.SetPort(c.Port())
	t.Log(p.Port())
	if err != nil {
		t.Error(err)
	}
	n, err := network.Connect(p)
	if err != nil {
		t.Error(err)
	}
	res, err := n.Ping()
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("node should be alive")
	}

	if p.Port() != port {
		t.Error("ports don't match")
	}
	network.Kill()
}

func TestGetClusterInfo(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", 0, &part)
	ver := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(ver, true)
	i := info.CreateInMemoryClusterInfo()
	peer1 := peer.CreateLocalPeer("testpeer1", 0, &part)
	peer2 := peer.CreateLocalPeer("testpeer2", 0, &part)
	i.Add(peer1, peerInfo)
	i.Add(peer2, peerInfo)

	network := CreateRpcNetwork()

	c, err := network.Serve(p, i)

	p.SetPort(c.Port())
	if err != nil {
		t.Error(err)
	}

	server := peer.CreateLocalPeer(p.Name(), p.Port(), p.Partition())

	n, err := network.Connect(server)
	if err != nil {
		t.Error(err)
	}
	resp, err := n.GetClusterInfo()
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 2 {
		t.Error("size doesn't match")
	}
	network.Kill()
}

func TestDeadNode(t *testing.T) {
	port := uint16(63447)
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", port, &part)
	// i := info.CreateInMemoryClusterInfo()
	// _, err := CreateRpcServer(p, i)

	// if err != nil {
	// 	t.Error(err)
	// }

	network := CreateRpcNetwork()
	if _, err := network.Connect(p); err == nil {
		t.Error("node is not dead")
	}
}
