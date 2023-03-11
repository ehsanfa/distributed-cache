package peer

import (
	"bytes"
	"dbcache/cluster/version"
	"encoding/gob"
)

type SimplePeerInfo struct {
	version version.Version
	isAlive bool
}

type marshalPeerInfo struct {
	Version []byte
	IsAlive bool
}

func CreateSimplePeerInfo(ver version.Version, isAlive bool) *SimplePeerInfo {
	return &SimplePeerInfo{version: ver, isAlive: isAlive}
}

func (i *SimplePeerInfo) MarkAsDead() {
	i.isAlive = false
}

func (i *SimplePeerInfo) Version() version.Version {
	return i.version
}

func (i *SimplePeerInfo) IsAlive() bool {
	return i.isAlive
}

func (i *SimplePeerInfo) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	mv, err := i.version.MarshalBinary()
	if err != nil {
		return make([]byte, 0), err
	}
	if err := enc.Encode(marshalPeerInfo{
		Version: mv,
		IsAlive: i.isAlive,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *SimplePeerInfo) UnmarshalBinary(data []byte) error {
	m := marshalPeerInfo{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	v := version.CreateGenClockVersion()
	if e := v.UnmarshalBinary(m.Version); e != nil {
		return e
	}

	i.version = v
	i.isAlive = m.IsAlive
	return nil
}
