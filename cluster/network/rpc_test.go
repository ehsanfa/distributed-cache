package network

import (
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"testing"
)

func TestWithRandomPort(t *testing.T) {
	p := peer.CreateLocalPeer("0.0.0.0", 0)
	p.SetPartition(partition.CreateSimplePartition("0"))
	i := info.CreateInMemoryClusterInfo()
	c, err := CreateRpcServer(p, i)
	p.SetPort(c.Port())
	t.Log(p.Port())
	if err != nil {
		t.Error(err)
	}
	res, err := c.Ping()
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("node should be alive")
	}
}

func TestWithSpecificPort(t *testing.T) {
	port := uint16(63447)
	p := peer.CreateLocalPeer("0.0.0.0", port)
	p.SetPartition(partition.CreateSimplePartition("0"))
	i := info.CreateInMemoryClusterInfo()
	c, err := CreateRpcServer(p, i)
	p.SetPort(c.Port())
	t.Log(p.Port())
	if err != nil {
		t.Error(err)
	}
	res, err := c.Ping()
	if err != nil {
		t.Error(err)
	}
	if !res {
		t.Error("node should be alive")
	}

	if p.Port() != port {
		t.Error("ports don't match")
	}
}

func TestGetClusterInfo(t *testing.T) {
	p := peer.CreateLocalPeer("0.0.0.0", 0)
	ver := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(ver, true)
	p.SetPartition(partition.CreateSimplePartition("0"))
	i := info.CreateInMemoryClusterInfo()
	peer1 := peer.CreateLocalPeer("testpeer1", 0)
	peer2 := peer.CreateLocalPeer("testpeer2", 0)
	i.Add(peer1, peerInfo)
	i.Add(peer2, peerInfo)
	c, err := CreateRpcServer(p, i)
	p.SetPort(c.Port())
	if err != nil {
		t.Error(err)
	}
	resp, err := c.GetClusterInfo()
	if err != nil {
		t.Error(err)
	}
	panic(resp)
	if len(resp) != 2 {
		t.Error("size doesn't match")
	}
}
