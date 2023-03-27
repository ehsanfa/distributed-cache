package peer

import (
	"bytes"
	rpcVersion "dbcache/cluster/network/rpc/types/version"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"encoding/gob"
)

type PeerInfo struct {
	Pi peer.PeerInfo
}

type marshalPeerInfo struct {
	Version []byte
	IsAlive bool
}

func (i *PeerInfo) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	v := rpcVersion.Version{Ver: i.Pi.Version()}
	mv, err := v.MarshalBinary()
	if err != nil {
		return make([]byte, 0), err
	}
	if err := enc.Encode(marshalPeerInfo{
		Version: mv,
		IsAlive: i.Pi.IsAlive(),
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *PeerInfo) UnmarshalBinary(data []byte) error {
	m := marshalPeerInfo{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	v := rpcVersion.Version{Ver: version.CreateGenClockVersion(0)}
	if e := v.UnmarshalBinary(m.Version); e != nil {
		return e
	}
	i.Pi = peer.CreateSimplePeerInfo(v.Ver, m.IsAlive)
	return nil
}
