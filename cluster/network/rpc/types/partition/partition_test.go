package partition

import (
	"bytes"
	"dbcache/cluster/partition"
	"encoding/gob"
	"testing"
)

func TestMarshal(t *testing.T) {
	p := Partition{partition.CreateSimplePartition("partition_1")}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		t.Error(err)
	}
	s := Partition{}
	reader := bytes.NewReader(buf.Bytes())
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&s); err != nil {
		t.Error(err)
	}

	if p.Part.Name() != s.Part.Name() {
		t.Error("partition names don't match")
	}
}
