package cache

import (
	"bytes"
	"dbcache/cluster/cacher"
	"encoding/gob"
)

type RpcCacheValue struct {
	Value cacher.CacheValue
}

type marshalCacheValue struct {
	Value   string
	Version int
}

func (v RpcCacheValue) GetCacheValue() cacher.CacheValue {
	return v.Value
}

func (v RpcCacheValue) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshalCacheValue{v.Value.GetValue(), v.Value.Version()}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *RpcCacheValue) UnmarshalBinary(data []byte) error {
	m := &marshalCacheValue{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return err
	}
	v.Value = cacher.NewVersionBasedCacheValue(m.Value, m.Version)
	return nil
}
