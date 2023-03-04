package info

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"testing"
)

func TestInfo(t *testing.T) {
	info := CreateInMemoryClusterInfo()
	part := partition.CreateSimplePartition("part1")
	version := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(version, true)
	p := peer.CreateLocalPeer("peer1", 12564)
	p.SetPartition(part)
	info.Add(p, peerInfo)

	for peer1, peerInfo1 := range info.All() {
		if !peer1.IsSameAs(peer.CreateLocalPeer("peer1", 12564)) {
			t.Error("comparing all method failed peer", peer1, peer.CreateLocalPeer("peer1", 12564))
		}
		if peerInfo1 != peerInfo {
			t.Error("comparing all method failed peerinfo", peerInfo1, peerInfo)
		}
	}

	for _, peer1 := range info.List() {
		if !peer1.IsSameAs(peer.CreateLocalPeer("peer1", 12564)) {
			t.Error("comparing all method failed peer", peer1, peer.CreateLocalPeer("peer1", 12564))
		}
	}
}

func TestIsPeerKnownAlive(t *testing.T) {
	info := CreateInMemoryClusterInfo()
	part := partition.CreateSimplePartition("part1")
	version := version.CreateGenClockVersion()
	peerInfo := peer.CreateSimplePeerInfo(version, true)
	p := peer.CreateLocalPeer("peer1", 12564)
	p.SetPartition(part)
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

func TestUpdate(t *testing.T) {
	info := CreateInMemoryClusterInfo()
	peerInfo1 := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(), true)
	p1 := peer.CreateLocalPeer("peer1", 12564)
	info.Add(p1, peerInfo1)
	peerInfo2 := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(), true)
	p2 := peer.CreateLocalPeer("peer2", 12564)
	info.Update(map[peer.Peer]peer.PeerInfo{p2: peerInfo2})

	if peerInfo1.Version().Number() != 0 || peerInfo2.Version().Number() != 0 {
		t.Error("initial versions don't match", peerInfo2.Version().Number())
	}
	newPeerInfo := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(), true)
	newPeerInfo.Version().Increment()
	newPeerInfo.Version().Increment()
	newPeerInfo.Version().Increment()
	info.Update(map[peer.Peer]peer.PeerInfo{p2: newPeerInfo})
	res, _ := info.Get(p2)
	if res.Version().Number() != newPeerInfo.Version().Number() {
		t.Error("updated versions don't match", res.Version().Number(), newPeerInfo.Version().Number())
	}

	outdatedPeerInfo := peer.CreateSimplePeerInfo(version.CreateGenClockVersion(), true)
	info.Update(map[peer.Peer]peer.PeerInfo{p2: outdatedPeerInfo})
	res, _ = info.Get(p2)
	if res.Version().Number() == outdatedPeerInfo.Version().Number() {
		t.Error("updated versions don't match", res.Version().Number(), outdatedPeerInfo.Version().Number())
	}
}
