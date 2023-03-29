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

func TestGetsInfoFromSeeder(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	cache := cacher.CreateInMemoryCache()
	buff := buffer.CreateInMemoryBuffer()
	seederPort := 45654
	go func() {
		seeder := peer.CreateLocalPeer("0.0.0.0", uint16(seederPort), &part)
		peer2 := peer.CreateLocalPeer("0.0.0.0", 45656, &part)
		peer2Info := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(1), true)
		peer3 := peer.CreateLocalPeer("0.0.0.0", 45657, &part)
		peer3Info := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(1), true)
		peer4 := peer.CreateLocalPeer("0.0.0.0", 45658, &part)
		peer4Info := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(1), true)

		seederInfo := info.CreateInMemoryClusterInfo()
		seederInfo.Add(peer2, peer2Info)
		seederInfo.Add(peer3, peer3Info)
		seederInfo.Add(peer4, peer4Info)
		_, err := rpc.CreateRpcNetwork(seeder, seederInfo, cache, buff)
		if err != nil {
			t.Error(err)
		}
	}()

	peer1 := peer.CreateLocalPeer("0.0.0.0", 45655, &part)

	info1 := info.CreateInMemoryClusterInfo()
	network1, err := rpc.CreateRpcNetwork(peer1, info1, cache, buff)
	if err != nil {
		t.Error(err)
	}

	seeder := peer.CreateLocalPeer("0.0.0.0", uint16(seederPort), &part)
	n, err := network1.Connect(seeder, 10*time.Second)
	if err != nil {
		t.Error(err)
	}
	resp, err := n.GetClusterInfo()
	if err != nil {
		t.Error(err)
	}
	if len(resp) != 3 {
		t.Error("size doesn't match")
	}

	gossip := CreateGossipNetwork(network1, info1, seeder)
	log.Printf("info in test: %p", info1)
	gossip.Start()
	log.Println("sleeping")
	time.Sleep(interval + 2*time.Second)
	if ok := info1.IsPeerKnown(peer.CreateLocalPeer("0.0.0.0", 45658, part)); !ok {
		t.Error("expected to receive cluster info")
	}
}
