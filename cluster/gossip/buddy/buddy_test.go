package buddy

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"testing"
)

func TestBuddy(t *testing.T) {
	part := partition.CreateSimplePartition("120")
	peer1 := peer.CreateLocalPeer("peer1", 12564)
	peer1.SetPartition(part)
	peer2 := peer.CreateLocalPeer("peer2", 12564)
	peer2.SetPartition(part)
	peer3 := peer.CreateLocalPeer("peer3", 12564)
	peer3.SetPartition(part)
	buddies := CreateInMemoryBuddies(2)
	if !buddies.CanAcceptBuddyRequest() {
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
	if buddies.CanAcceptBuddyRequest() {
		t.Error("can accept buddies failed")
	}

	all := buddies.AllBuddies()
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
	if !buddies.IsBuddyWith(peer3) {
		t.Error("isbuddywith method failed")
	}
}
