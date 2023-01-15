package info

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"reflect"
	"testing"
)

func TestInfo(t *testing.T) {
	info := CreateInMemoryClusterInfo()
	part := partition.CreateSimplePartition("part1")
	version := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(version, true)
	p := peer.CreateLocalPeer("peer1", 12564, part)
	info.Add(p, peerInfo)

	for peer1, peerInfo1 := range info.All() {
		if peer1 != peer.CreateLocalPeer("peer1", 12564, part) {
			t.Error("comparing all method failed peer", peer1, peer.CreateLocalPeer("peer1", 12564, part))
		}
		if peerInfo1 != peerInfo {
			t.Error("comparing all method failed peerinfo", peerInfo1, peerInfo)
		}
	}

	if !reflect.DeepEqual(info.List(), []peer.Peer{peer.CreateLocalPeer("peer1", 12564, part)}) {
		t.Error("comparing list method failed", info.List(), []peer.Peer{peer.CreateLocalPeer("peer1", 12564, part)})
	}
}

func TestIsPeerKnownAlive(t *testing.T) {
	info := CreateInMemoryClusterInfo()
	part := partition.CreateSimplePartition("part1")
	version := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(version, true)
	p := peer.CreateLocalPeer("peer1", 12564, part)
	info.Add(p, peerInfo)

	if !info.IsPeerKnown(p) {
		t.Error("isPeerKnown failed")
	}

	if !info.IsPeerAlive(p) {
		t.Error("isPeerAlive failed")
	}

	info.Remove(p)

	if info.IsPeerKnown(p) {
		t.Error("isPeerKnown failed")
	}

	if info.IsPeerAlive(p) {
		t.Error("isPeerAlive failed")
	}
}
