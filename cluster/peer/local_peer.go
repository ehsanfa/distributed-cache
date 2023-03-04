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

func CreateLocalPeer(name string, port uint16) Peer {
	return &LocalPeer{name: name, port: port}
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
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *LocalPeer) UnmarshalBinary(data []byte) error {
	l := &LocalPeer{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&l); err != nil {
		return err
	}
	p.name = l.name
	p.port = l.port
	return nil
}
