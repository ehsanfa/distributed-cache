package rpc

import (
	"bytes"
	"dbcache/cluster/cacher"
	rpcCache "dbcache/cluster/network/rpc/types/cache"
	"encoding/gob"
)

type GetCacheResponse struct {
	Cache map[string]cacher.CacheValue
}

type marshalCacheResponse struct {
	Resp []marshalCache
}

type marshalCache struct {
	Key   string
	Value []byte
}

func (n *RpcNode) GetCache() (map[string]cacher.CacheValue, error) {
	resp := new(GetCacheResponse)
	err := n.client.Call(n.rpcAction("RpcGetCache"), struct{}{}, &resp)
	return resp.Cache, err
}

func (n *RpcNode) RpcGetCache(p struct{}, resp *GetCacheResponse) error {
	*resp = GetCacheResponse{hostNetwork.cache.All()}
	return nil
}

func (r *GetCacheResponse) MarshalBinary() (data []byte, err error) {
	var c []marshalCache
	for k, v := range r.Cache {
		rpccv := rpcCache.RpcCacheValue{Value: v}
		mv, err := rpccv.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		mc := marshalCache{Key: k, Value: mv}
		c = append(c, mc)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshalCacheResponse{
		Resp: c,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *GetCacheResponse) UnmarshalBinary(data []byte) error {
	c := make(map[string]cacher.CacheValue)
	mcr := &marshalCacheResponse{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&mcr); err != nil {
		return err
	}
	cacheValue := rpcCache.RpcCacheValue{Value: cacher.NewVersionBasedCacheValue("", 0)}
	for _, mc := range mcr.Resp {
		if e := cacheValue.UnmarshalBinary(mc.Value); e != nil {
			return e
		}
		cv := rpcCache.RpcCacheValue{Value: cacheValue.Value}
		c[mc.Key] = cv.Value
	}
	r.Cache = c
	return nil
}
