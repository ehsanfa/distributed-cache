package node

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/gossip/buddy"
	"dbcache/cluster/info"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"net/rpc"
	"testing"
)

type mockCluster struct {
}

func createMockCluster() mockCluster {
	return mockCluster{}
}

func createCacheValue(value string) cacher.CacheValue {
	v := cacher.NewCacheValue()
	v.Value = value
	return v
}

func (m mockCluster) GetClusterInfo(p peer.Peer) (map[peer.Peer]peer.PeerInfo, error) {
	v := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(v, true)
	info := map[peer.Peer]peer.PeerInfo{p: peerInfo}
	return info, nil
}

func (m mockCluster) AskForParition(peer.Peer) (partition.Partition, error) {
	return partition.CreateSimplePartition("0"), nil
}

func (m mockCluster) GetCache(peer.Peer) (map[string]cacher.CacheValue, error) {
	v := createCacheValue("hooshang")
	return map[string]cacher.CacheValue{"hasan": v}, nil
}

func createNode(p peer.Peer, isSeeder bool) *Node {
	return &Node{
		peer:     p,
		isSeeder: isSeeder,
		// partition:   partition,
		info:        info.CreateInMemoryClusterInfo(),
		cache:       cacher.CreateInMemoryCache(),
		buffer:      buffer.CreateInMemoryBuffer(),
		connections: make(map[peer.Peer]*rpc.Client),
		buddies:     buddy.CreateInMemoryBuddies(maxBuddyNum),
		cluster:     createMockCluster(),
	}
}

func TestIntroduce(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	s := peer.CreateLocalPeer("seeder", 12345)
	s.SetParition(part)
	p1 := peer.CreateLocalPeer("random_peer1", 12345)
	p1.SetParition(part)
	p := peer.CreateLocalPeer("node", 12345)
	p.SetParition(part)
	n := createNode(p, false)
	n.SetSeeder(s)
	err := n.Introduce()
	if err != nil {
		t.Error(err)
	}

	if !n.info.IsPeerKnown(s) {
		t.Error("seeder should be known")
	}
	if n.info.IsPeerKnown(p1) {
		t.Error("random peer shouldn't be known")
	}

	if !n.cache.Exists("hasan") {
		t.Error("cache doesn't exist")
	}
	v, _ := n.cache.Get("hasan")
	if v != createCacheValue("hooshang") {
		t.Error("cache values don't match")
	}
}
