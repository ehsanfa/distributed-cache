package partition

import (
	"bytes"
	"dbcache/cluster/partition"
	"encoding/gob"
)

type Partition struct {
	Part partition.Partition
}

type marshal struct {
	Name string
}

func (p Partition) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshal{
		Name: p.Part.Name(),
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Partition) UnmarshalBinary(data []byte) error {
	m := &marshal{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	p.Part = partition.CreateSimplePartition(m.Name)
	return nil
}
