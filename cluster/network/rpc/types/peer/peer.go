package peer

import (
	"bytes"
	rpcPartition "dbcache/cluster/network/rpc/types/partition"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"encoding/gob"
)

type Peer struct {
	Peer peer.Peer
}

type marshalPeer struct {
	Name      string
	Port      uint16
	Partition []byte
}

func (p *Peer) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	mp := marshalPeer{
		Name: p.Peer.Name(),
		Port: p.Peer.Port(),
	}
	var mpp []byte
	if p.Peer.Partition() != nil {
		part := rpcPartition.Partition{Part: p.Peer.Partition()}
		mpp, err = part.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		mp.Partition = mpp
	}

	if err := enc.Encode(mp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Peer) UnmarshalBinary(data []byte) error {
	m := &marshalPeer{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	p.Peer = peer.CreateLocalPeer(m.Name, m.Port)
	if m.Partition != nil {
		mp := rpcPartition.Partition{Part: partition.CreateSimplePartition("")}
		if e := mp.UnmarshalBinary(m.Partition); e != nil {
			return e
		}
		p.Peer.SetPartition(mp.Part)
	}
	return nil
}
