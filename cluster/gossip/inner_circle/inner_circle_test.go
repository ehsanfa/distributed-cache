package innercircle

import (
	"dbcache/cluster/peer"
	"testing"
)

func TestInnerCircle(t *testing.T) {
	self := peer.CreateLocalPeer("self", 12564)
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	buddies := CreateInMemoryBuddies(2, self)
	if buddies.canAdd(self) {
		t.Error("cannot add self as buddy")
	}
	if !buddies.canAcceptBuddyRequest() {
		t.Error("can accept buddies failed")
	}
	if !buddies.Add(peer1) {
		t.Error("adding a new buddy failed")
	}
	if !buddies.Add(peer2) {
		t.Error("adding a new buddy failed")
	}
	if buddies.Add(peer3) {
		t.Error("adding a new buddy failed")
	}
	if buddies.canAcceptBuddyRequest() {
		t.Error("can accept buddies failed")
	}

	all := buddies.All()
	toCompare := map[peer.Peer]bool{peer1: true, peer2: true}
	for peer := range all {
		if _, ok := toCompare[peer]; !ok {
			t.Error("all method failed", all, toCompare)
		}
	}

	if buddies.Count() != 2 {
		t.Error("count method failed")
	}

	buddies.Remove(peer1)
	if buddies.Count() != 1 {
		t.Error("count method failed")
	}

	buddies.Remove(peer2)
	if !buddies.IsEmpty() {
		t.Error("IsEmpty method failed")
	}

	buddies.Add(peer3)
	if !buddies.isBuddyWith(peer3) {
		t.Error("isbuddywith method failed")
	}
}

func TestReplace(t *testing.T) {
	self := peer.CreateLocalPeer("self", 12564)
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	buddies := CreateInMemoryBuddies(3, self)
	buddies.Add(peer1)
	buddies.Add(peer2)
	buddies.Add(peer3)
	peer4 := peer.CreateLocalPeer("peer4", 12564)
	buddies.Replace(peer1, peer4)
	if buddies.isBuddyWith(peer1) {
		t.Error("expected not to be buddy with", peer1)
	}

	if !buddies.isBuddyWith(peer4) {
		t.Error("expected to be buddy with", peer4)
	}
}

func TestShuffle(t *testing.T) {
	self := peer.CreateLocalPeer("self", 12564)
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	buddies := CreateInMemoryBuddies(3, self)
	buddies.Add(peer1)
	buddies.Add(peer2)
	buddies.Add(peer3)

	peer4 := peer.CreateLocalPeer("peer4", 12564)
	peer5 := peer.CreateLocalPeer("peer5", 12564)
	peer6 := peer.CreateLocalPeer("peer6", 12564)
	peer7 := peer.CreateLocalPeer("peer7", 12564)
	peer8 := peer.CreateLocalPeer("peer8", 12564)
	info := []peer.Peer{peer1, peer2, peer3, peer4, peer5, peer6, peer7, peer8}
	buddies.Shuffle(info)
	if buddies.isBuddyWith(peer1) {
		t.Error("expected not to be buddy with peers", peer1)
	}

	if buddies.isBuddyWith(peer2) {
		t.Error("expected not to be buddy with peers", peer2)
	}

	if buddies.isBuddyWith(peer3) {
		t.Error("expected not to be buddy with peers", peer3)
	}

	if buddies.Count() != 3 {
		t.Errorf("expected to have %d buddies", 3)
	}
}

func TestShuffleWhenMaxNumIsNotMet(t *testing.T) {
	self := peer.CreateLocalPeer("self", 12564)
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	buddies := CreateInMemoryBuddies(5, self)
	buddies.Add(peer1)
	buddies.Add(peer2)
	buddies.Add(peer3)
	peer4 := peer.CreateLocalPeer("peer4", 12564)
	info := []peer.Peer{peer1, peer2, peer3, peer4}
	buddies.Shuffle(info)
	if buddies.Count() != 4 {
		t.Errorf("expected to have %d buddies but has %d", 4, buddies.Count())
	}
	if !buddies.isBuddyWith(peer4) {
		t.Errorf("expected to be buddy with %s", peer4)
	}
}

func TestShuffleWhenInfoNotEmptyButBuddiesIs(t *testing.T) {
	self := peer.CreateLocalPeer("self", 12564)
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	buddies := CreateInMemoryBuddies(5, self)
	info := []peer.Peer{peer1, peer2, peer3}
	buddies.Shuffle(info)
	if buddies.Count() != 3 {
		t.Errorf("expected to have %d buddies but has %d", 3, buddies.Count())
	}
}
