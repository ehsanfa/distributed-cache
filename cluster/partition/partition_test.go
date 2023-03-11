package partition

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestMarshal(t *testing.T) {
	p := CreateSimplePartition("partition_1")
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		t.Error(err)
	}
	s := SimplePartition{}
	reader := bytes.NewReader(buf.Bytes())
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&s); err != nil {
		t.Error(err)
	}

	if p.Name() != s.Name() {
		t.Error("partition names don't match")
	}
}
