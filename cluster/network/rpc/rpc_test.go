package rpc

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"testing"
	"time"
)

func TestWithRandomPort(t *testing.T) {
	p := peer.CreateLocalPeer("0.0.0.0", 0)
	i := info.CreateInMemoryClusterInfo()
	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	network, err := CreateRpcNetwork(p, i, cache, buff)
	if err != nil {
		t.Error(err)
	}
	p = network.Peer()

	n, err := network.Connect(p, 10*time.Second)
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
	// t.Error(p.Port())
	t.Cleanup(network.Kill)
}

func TestGetClusterInfo(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	go func() {
		p := peer.CreateLocalPeer("0.0.0.0", 5666)
		ver := version.CreateGenClockVersion(0)
		peerInfo := peer.CreateSimplePeerInfo(peer.Cacher, ver, true)
		i := info.CreateInMemoryClusterInfo()
		peer1 := peer.CreateLocalPeer("testpeer1", 0)
		peer1 = peer1.SetPartition(part)
		peer2 := peer.CreateLocalPeer("testpeer2", 0)
		peer2 = peer2.SetPartition(part)
		i.Add(peer1, peerInfo)
		i.Add(peer2, peerInfo)

		cache1 := cacher.CreateInMemoryCache()
		buff1 := buffer.CreateInMemoryBuffer()

		_, err := CreateRpcNetwork(p, i, cache1, buff1)
		if err != nil {
			t.Error(err)
		}
	}()

	p2 := peer.CreateLocalPeer("0.0.0.0", 4666)
	i2 := info.CreateInMemoryClusterInfo()
	cache2 := cacher.CreateInMemoryCache()
	buff2 := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p2, i2, cache2, buff2)

	server := peer.CreateLocalPeer("0.0.0.0", 5666)

	n, err := network.Connect(server, 10*time.Second)
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
	// t.Cleanup(network.Kill)
}

func TestGetClusterInfoTwo(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", 0)
	ver := version.CreateGenClockVersion(0)
	peerInfo := peer.CreateSimplePeerInfo(peer.Cacher, ver, true)
	i := info.CreateInMemoryClusterInfo()
	peer1 := peer.CreateLocalPeer("testpeer1", 0)
	peer1 = peer1.SetPartition(part)
	peer2 := peer.CreateLocalPeer("testpeer2", 0)
	peer2 = peer2.SetPartition(part)
	i.Add(peer1, peerInfo)
	i.Add(peer2, peerInfo)

	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()

	network, err := CreateRpcNetwork(p, i, cache, buff)
	if err != nil {
		t.Error(err)
	}
	p = network.Peer()

	server := peer.CreateLocalPeer(p.Name(), p.Port())

	n, err := network.Connect(server, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	resp, err := n.GetClusterInfo()
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 2 {
		t.Error("size doesn't match", p.Port())
	}
	t.Cleanup(network.Kill)
}

func TestWithSpecificPort(t *testing.T) {
	port := uint16(63447)
	p := peer.CreateLocalPeer("0.0.0.0", port)
	i := info.CreateInMemoryClusterInfo()
	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	network, err := CreateRpcNetwork(p, i, cache, buff)
	t.Log(p.Port())
	if err != nil {
		t.Error(err)
	}
	n, err := network.Connect(p, 10*time.Second)
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
	// t.Error(p.Port())
	t.Cleanup(network.Kill)
}

func TestGetAllCache(t *testing.T) {
	go func() {
		p := peer.CreateLocalPeer("0.0.0.0", 5667)
		ver := version.CreateGenClockVersion(0)
		peerInfo := peer.CreateSimplePeerInfo(peer.Cacher, ver, true)
		i := info.CreateInMemoryClusterInfo()
		peer1 := peer.CreateLocalPeer("testpeer1", 0)
		peer2 := peer.CreateLocalPeer("testpeer2", 0)
		i.Add(peer1, peerInfo)
		i.Add(peer2, peerInfo)

		cache1 := cacher.CreateInMemoryCache()
		cache1.Set("hasan", cacher.NewVersionBasedCacheValue("hooshang", 1))

		buff := buffer.CreateInMemoryBuffer()

		_, err := CreateRpcNetwork(p, i, cache1, buff)
		if err != nil {
			t.Error(err)
		}
	}()

	p2 := peer.CreateLocalPeer("0.0.0.0", 4667)
	i2 := info.CreateInMemoryClusterInfo()
	cache2 := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p2, i2, cache2, buff)

	server := peer.CreateLocalPeer("0.0.0.0", 5667)

	n, err := network.Connect(server, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	resp, err := n.GetCache()
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 1 {
		t.Error("size doesn't match")
	}
	if resp["hasan"].GetValue() != "hooshang" {
		t.Error("wrong value provided")
	}
	if resp["hasan2"] != nil {
		t.Error("wrong value provided")
	}
	// t.Cleanup(network.Kill)
}

func TestGet(t *testing.T) {
	p1 := peer.CreateLocalPeer("0.0.0.0", 45655)
	cache1 := cacher.CreateInMemoryCache()
	cache1.Set("hasan", cacher.NewVersionBasedCacheValue("hooshang", 1))
	info1 := info.CreateInMemoryClusterInfo()
	buffer1 := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p1, info1, cache1, buffer1)

	n, err := network.Connect(peer.CreateLocalPeer("0.0.0.0", 45655), 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	resp, _ := n.Get("hasan")
	if resp.GetValue() != "hooshang" {
		t.Error("expected to get correct cached value")
	}
}

func TestPut(t *testing.T) {
	p1 := peer.CreateLocalPeer("0.0.0.0", 45655)
	cache1 := cacher.CreateInMemoryCache()
	info1 := info.CreateInMemoryClusterInfo()
	buffer1 := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p1, info1, cache1, buffer1)
	n, err := network.Connect(peer.CreateLocalPeer("0.0.0.0", 45655), 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	val := cacher.NewVersionBasedCacheValue("hooshang", 3)
	err = n.Set("hasan", val)
	if err != nil {
		t.Error(err)
	}

	if v, _ := cache1.Get("hasan"); v.GetValue() != "hooshang" {
		t.Error("expected to get the right value after set")
	}
	if v, _ := cache1.Get("hasan"); v.Version() != 3 {
		t.Error("expected to get the right version after set")
	}
}

func TestIntroduction(t *testing.T) {
	go func() {
		p := peer.CreateLocalPeer("0.0.0.0", 5666)
		i := info.CreateInMemoryClusterInfo()

		cache1 := cacher.CreateInMemoryCache()
		buff1 := buffer.CreateInMemoryBuffer()

		_, err := CreateRpcNetwork(p, i, cache1, buff1)
		if err != nil {
			t.Error(err)
		}
	}()

	p2 := peer.CreateLocalPeer("0.0.0.0", 4666)
	i2 := info.CreateInMemoryClusterInfo()
	cache2 := cacher.CreateInMemoryCache()
	buff2 := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p2, i2, cache2, buff2)

	server := peer.CreateLocalPeer("0.0.0.0", 5666)

	n, err := network.Connect(server, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	err = n.Introduce(peer.Cacher, p2)
	if err != nil {
		t.Error(err)
	}

	resp, err := n.GetClusterInfo()
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 1 {
		t.Error("size doesn't match")
	}
}

func TestUpdateBuffer(t *testing.T) {

}
