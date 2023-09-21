package gossip

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/network/rpc"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"log"
	"testing"
	"time"
)

func TestGetsInfoFromSeederWhenInitialized(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	seederPort := 45654
	go func() {
		seeder := peer.CreateLocalPeer("0.0.0.0", uint16(seederPort))
		peer2 := peer.CreateLocalPeer("0.0.0.0", 45656)
		peer2 = peer2.SetPartition(part)
		peer2Info := peer.CreateSimplePeerInfo(peer.Cacher, version.CreateGenClockVersion(1), true)
		peer3 := peer.CreateLocalPeer("0.0.0.0", 45657)
		peer3 = peer3.SetPartition(part)
		peer3Info := peer.CreateSimplePeerInfo(peer.Cacher, version.CreateGenClockVersion(1), true)
		peer4 := peer.CreateLocalPeer("0.0.0.0", 45658)
		peer4 = peer4.SetPartition(part)
		peer4Info := peer.CreateSimplePeerInfo(peer.Cacher, version.CreateGenClockVersion(1), true)

		seederInfo := info.CreateInMemoryClusterInfo()
		seederInfo.Add(peer2, peer2Info)
		seederInfo.Add(peer3, peer3Info)
		seederInfo.Add(peer4, peer4Info)
		_, err := rpc.CreateRpcNetwork(seeder, seederInfo, cache, buff)
		if err != nil {
			t.Error(err)
		}
	}()

	peer1 := peer.CreateLocalPeer("0.0.0.0", 45655)

	info1 := info.CreateInMemoryClusterInfo()
	network1, err := rpc.CreateRpcNetwork(peer1, info1, cache, buff)
	if err != nil {
		t.Error(err)
	}

	seeder := peer.CreateLocalPeer("0.0.0.0", uint16(seederPort))

	gossip, err := CreateGossipNetwork(network1, info1, seeder, false, peer1)
	if err != nil {
		t.Error(err)
	}
	log.Printf("info in test: %p", info1)
	gossip.Start()
	log.Println("sleeping")
	time.Sleep(gossipInterval + 2*time.Second)
	if ok := info1.IsPeerKnown(peer.CreateLocalPeer("0.0.0.0", 45658)); !ok {
		t.Error("expected to receive cluster info")
	}
}

func TestFailsWhenNoSeeder(t *testing.T) {
	peer1 := peer.CreateLocalPeer("0.0.0.0", 45659)
	info1 := info.CreateInMemoryClusterInfo()
	cache1 := cacher.CreateInMemoryCache()
	buff1 := buffer.CreateInMemoryBuffer()
	network1, err := rpc.CreateRpcNetwork(peer1, info1, cache1, buff1)
	if err != nil {
		t.Error(err)
	}
	if _, err := CreateGossipNetwork(network1, info1, nil, false, peer1); err == nil {
		t.Error("expected to see error")
	}
}

func TestSeederStandalone(t *testing.T) {
	peer1 := peer.CreateLocalPeer("0.0.0.0", 45659)
	info1 := info.CreateInMemoryClusterInfo()
	cache1 := cacher.CreateInMemoryCache()
	buff1 := buffer.CreateInMemoryBuffer()
	network1, err := rpc.CreateRpcNetwork(peer1, info1, cache1, buff1)
	if err != nil {
		t.Error(err)
	}
	gossip, err := CreateGossipNetwork(network1, info1, nil, true, peer1)
	if err != nil {
		t.Error(err)
	} else {
		gossip.Start()
	}

	time.Sleep(5 * time.Second)
}

func TestInfoIsUpdatedWhenNodeIsDead(t *testing.T) {
	// TODO
}
