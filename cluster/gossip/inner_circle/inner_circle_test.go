package innercircle

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"testing"
)

func TestInnerCircle(t *testing.T) {
	part := partition.CreateSimplePartition("120")
	peer1 := peer.CreateLocalPeer("peer1", 12564, part)
	peer2 := peer.CreateLocalPeer("peer2", 12564, part)
	peer3 := peer.CreateLocalPeer("peer3", 12564, part)
	buddies := CreateInMemoryBuddies(2)
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
