package peer

import (
	"bytes"
	"dbcache/cluster/partition"
	"encoding/gob"
)

type LocalPeer struct {
	name      string
	port      uint16
	partition partition.Partition
}

type marshalPeer struct {
	Name      string
	Port      uint16
	Partition []byte
}

func CreateLocalPeer(name string, port uint16, part partition.Partition) Peer {
	return &LocalPeer{name: name, port: port, partition: part}
}

func (p *LocalPeer) GetName() string {
	return p.name
}

func (p *LocalPeer) Name() string {
	return p.name
}

func (p *LocalPeer) GePort() uint16 {
	return p.port
}

func (p *LocalPeer) Port() uint16 {
	return p.port
}

func (p *LocalPeer) GetPartition() partition.Partition {
	return p.partition
}

func (p *LocalPeer) Partition() partition.Partition {
	return p.partition
}

func (p *LocalPeer) IsSameAs(other Peer) bool {
	return p.Name() == other.Name() && p.Port() == other.Port()
}

func (p *LocalPeer) SetPartition(part partition.Partition) {
	p.partition = part
}

func (p *LocalPeer) SetPort(port uint16) {
	p.port = port
}

func (p *LocalPeer) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if p.partition == nil {
		panic("Empty partition")
	}
	mp, err := p.partition.MarshalBinary()
	if err != nil {
		return make([]byte, 0), err
	}
	if err := enc.Encode(marshalPeer{
		Name:      p.name,
		Port:      p.port,
		Partition: mp,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *LocalPeer) UnmarshalBinary(data []byte) error {
	m := &marshalPeer{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	mp := partition.CreateSimplePartition("")
	if e := mp.UnmarshalBinary(m.Partition); e != nil {
		return e
	}
	p.name = m.Name
	p.port = m.Port
	p.partition = &mp
	return nil
}
