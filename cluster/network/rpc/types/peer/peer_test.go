package peer

import (
	"bytes"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"encoding/gob"
	"testing"
)

func TestMarshalPeerInfo(t *testing.T) {
	ver := version.CreateGenClockVersion(0)
	ver = ver.Increment().(version.GenClock)
	ver = ver.Increment().(version.GenClock)
	ver = ver.Increment().(version.GenClock)
	pi := PeerInfo{Pi: peer.CreateSimplePeerInfo(ver, true)}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&pi); err != nil {
		t.Error(err)
	}
	i := PeerInfo{}
	reader := bytes.NewReader(buf.Bytes())
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&i); err != nil {
		t.Error(err)
	}

	if pi.Pi.IsAlive() != i.Pi.IsAlive() {
		t.Error("peer isAlive don't match")
	}

	if i.Pi.Version().Number() != 3 {
		t.Error("version number is incorrect")
	}

	if pi.Pi.Version().Number() != i.Pi.Version().Number() {
		t.Error("version numbers don't match")
	}
}
