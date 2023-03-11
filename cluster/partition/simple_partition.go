package partition

import (
	"bytes"
	"encoding/gob"
)

type SimplePartition struct {
	name string
}

type marshal struct {
	Name string
}

func (p SimplePartition) Name() string {
	return p.name
}

func (p SimplePartition) IsEmpty() bool {
	return p == SimplePartition{}
}

func CreateSimplePartition(name string) SimplePartition {
	return SimplePartition{name: name}
}

func (p SimplePartition) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshal{
		Name: p.name,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *SimplePartition) UnmarshalBinary(data []byte) error {
	m := &marshal{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	p.name = m.Name
	return nil
}
