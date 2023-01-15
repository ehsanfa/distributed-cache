package peer

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/version"
	"testing"
)

func TestPeer(t *testing.T) {
	part := partition.CreateSimplePartition("120")
	p := CreateLocalPeer("peer1", 12564, part)
	if p.Name() != "peer1" {
		t.Error("name doesn't match")
	}

	if p.Port() != 12564 {
		t.Error("port doesn't match")
	}

	if p.Partition() != part {
		t.Error("partition doesn't match")
	}

	if !p.IsSameAs(CreateLocalPeer("peer1", 12564, part)) {
		t.Error("should be equal to other partition")
	}
}

func TestPeerInfo(t *testing.T) {
	ver := version.CreateGenClockVersion()
	info := CreateSimplePeerInfo(ver, true)
	if info.Version() != ver {
		t.Error("versions don't match", info.Version(), ver)
	}

	if info.IsAlive() != true {
		t.Error("isAlive doesn't match")
	}
}
