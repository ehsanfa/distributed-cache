package rpc

import (
	"dbcache/cluster/cacher"
)

type GetResponse struct {
	Key   string
	Value RpcCacheValue
	Ok    bool
}

type GetRequest struct {
	Key string
}

func (n *RpcNode) Get(key string) (cacher.CacheValue, error) {
	r := new(GetResponse)
	err := n.client.Call(n.rpcAction("RpcGet"), GetRequest{Key: key}, &r)
	return r.Value.value, err
}

func (n *RpcNode) RpcGet(r GetRequest, resp *GetResponse) error {
	resp.Key = r.Key
	var val cacher.CacheValue
	val, resp.Ok = hostNetwork.cache.Get(r.Key)
	resp.Value = RpcCacheValue{val}
	return nil
}

type SetResponse struct {
}

type SetRequest struct {
	Key   string
	Value *RpcCacheValue
}

func (n *RpcNode) Set(key string, value cacher.CacheValue) error {
	r := new(SetResponse)
	rv := &RpcCacheValue{value: value}
	err := n.client.Call(n.rpcAction("RpcSet"), SetRequest{Key: key, Value: rv}, &r)
	return err
}

func (n *RpcNode) RpcSet(r SetRequest, resp *SetResponse) error {
	err := n.network.cache.Set(r.Key, r.Value.value)
	return err
}
