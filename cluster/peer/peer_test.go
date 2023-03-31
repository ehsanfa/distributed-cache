package peer

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/version"
	"testing"
)

func TestPeer(t *testing.T) {
	part := partition.CreateSimplePartition("120")
	p := CreateLocalPeer("peer1", 12564)
	if p.Name() != "peer1" {
		t.Error("name doesn't match")
	}

	if p.Port() != 12564 {
		t.Error("port doesn't match")
	}

	p = p.SetPartition(&part)

	if p.Partition() != &part {
		t.Error("partition doesn't match")
	}

	if !p.IsSameAs(CreateLocalPeer("peer1", 12564)) {
		t.Error("should be equal to other peer")
	}
}

func TestPeerInfo(t *testing.T) {
	ver := version.CreateGenClockVersion(0)
	info := CreateSimplePeerInfo(ver, true)
	if info.Version() != ver {
		t.Error("versions don't match", info.Version(), ver)
	}

	if info.IsAlive() != true {
		t.Error("isAlive doesn't match")
	}
}

// func TestMarshalPeer(t *testing.T) {
// 	part := partition.CreateSimplePartition("partition_1")
// 	p := CreateLocalPeer("test1", 22, &part)
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	if err := enc.Encode(p); err != nil {
// 		t.Error(err)
// 	}
// 	l := LocalPeer{}
// 	reader := bytes.NewReader(buf.Bytes())
// 	dec := gob.NewDecoder(reader)
// 	if err := dec.Decode(&l); err != nil {
// 		t.Error(err)
// 	}

// 	if p.Name() != l.Name() {
// 		t.Error("peer names don't match")
// 	}

// 	if l.Partition().Name() != "partition_1" {
// 		t.Error("incorrect partition name")
// 	}

// 	if p.Partition().Name() != l.Partition().Name() {
// 		t.Error("partitions don't match")
// 	}
// }

// func TestMarshalPeerInfo(t *testing.T) {
// 	ver := version.CreateGenClockVersion(0)
// 	ver.Increment()
// 	ver.Increment()
// 	ver.Increment()
// 	pi := CreateSimplePeerInfo(ver, true)
// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	if err := enc.Encode(pi); err != nil {
// 		t.Error(err)
// 	}
// 	i := SimplePeerInfo{}
// 	reader := bytes.NewReader(buf.Bytes())
// 	dec := gob.NewDecoder(reader)
// 	if err := dec.Decode(&i); err != nil {
// 		t.Error(err)
// 	}

// 	if pi.IsAlive() != i.IsAlive() {
// 		t.Error("peer isAlive don't match")
// 	}

// 	if i.Version().Number() != 3 {
// 		t.Error("version number is incorrect")
// 	}

// 	if pi.Version().Number() != i.Version().Number() {
// 		t.Error("version numbers don't match")
// 	}
// }
