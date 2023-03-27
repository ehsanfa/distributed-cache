package rpc

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
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
	cache := cacheProvider{cacher.CreateInMemoryCache()}
	buff := buffer.CreateInMemoryBuffer()
	network, err := CreateRpcNetwork(p, i, &cache, buff)
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
	// t.Error(p.Port())
	t.Cleanup(network.Kill)
}

func TestGetClusterInfo(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	go func() {
		p := peer.CreateLocalPeer("0.0.0.0", 5666, &part)
		ver := version.CreateGenClockVersion(0)
		peerInfo := peer.CreateSimplePeerInfo(ver, true)
		i := info.CreateInMemoryClusterInfo()
		peer1 := peer.CreateLocalPeer("testpeer1", 0, &part)
		peer2 := peer.CreateLocalPeer("testpeer2", 0, &part)
		i.Add(peer1, peerInfo)
		i.Add(peer2, peerInfo)

		cache1 := cacher.CreateInMemoryCache()
		buff1 := buffer.CreateInMemoryBuffer()

		_, err := CreateRpcNetwork(p, i, &cacheProvider{cache1}, buff1)
		if err != nil {
			t.Error(err)
		}
	}()

	p2 := peer.CreateLocalPeer("0.0.0.0", 4666, &part)
	i2 := info.CreateInMemoryClusterInfo()
	cache2 := cacher.CreateInMemoryCache()
	buff2 := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p2, i2, &cacheProvider{cache2}, buff2)

	server := peer.CreateLocalPeer("0.0.0.0", 5666, &part)

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
	// t.Cleanup(network.Kill)
}

func TestGetClusterInfoTwo(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", 0, &part)
	ver := version.CreateGenClockVersion(0)
	peerInfo := peer.CreateSimplePeerInfo(ver, true)
	i := info.CreateInMemoryClusterInfo()
	peer1 := peer.CreateLocalPeer("testpeer1", 0, &part)
	peer2 := peer.CreateLocalPeer("testpeer2", 0, &part)
	i.Add(peer1, peerInfo)
	i.Add(peer2, peerInfo)

	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()

	network, err := CreateRpcNetwork(p, i, &cacheProvider{cache}, buff)
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
		t.Error("size doesn't match", p.Port())
	}
	t.Cleanup(network.Kill)
}

func TestWithSpecificPort(t *testing.T) {
	port := uint16(63447)
	part := partition.CreateSimplePartition("0")
	p := peer.CreateLocalPeer("0.0.0.0", port, &part)
	i := info.CreateInMemoryClusterInfo()
	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	network, err := CreateRpcNetwork(p, i, &cacheProvider{cache}, buff)
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
	// t.Error(p.Port())
	t.Cleanup(network.Kill)
}

type cacheProvider struct {
	cache cacher.Cache
}

func (c *cacheProvider) GetCache() map[string]cacher.CacheValue {
	t := make(map[string]cacher.CacheValue)
	for k, v := range c.cache.All() {
		t[k] = v
	}
	return t
}

func TestGetCache(t *testing.T) {

	part := partition.CreateSimplePartition("0")
	go func() {
		p := peer.CreateLocalPeer("0.0.0.0", 5667, &part)
		ver := version.CreateGenClockVersion(0)
		peerInfo := peer.CreateSimplePeerInfo(ver, true)
		i := info.CreateInMemoryClusterInfo()
		peer1 := peer.CreateLocalPeer("testpeer1", 0, &part)
		peer2 := peer.CreateLocalPeer("testpeer2", 0, &part)
		i.Add(peer1, peerInfo)
		i.Add(peer2, peerInfo)

		cache1 := cacher.CreateInMemoryCache()
		cache1.Set("hasan", cacher.NewVersionBasedCacheValue("hooshang", 1))

		buff := buffer.CreateInMemoryBuffer()

		_, err := CreateRpcNetwork(p, i, &cacheProvider{cache1}, buff)
		if err != nil {
			t.Error(err)
		}
	}()

	p2 := peer.CreateLocalPeer("0.0.0.0", 4667, &part)
	i2 := info.CreateInMemoryClusterInfo()
	cache2 := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	network, _ := CreateRpcNetwork(p2, i2, &cacheProvider{cache2}, buff)

	server := peer.CreateLocalPeer("0.0.0.0", 5667, &part)

	n, err := network.Connect(server)
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
