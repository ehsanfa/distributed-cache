package rpc

import (
	"dbcache/cluster/cacher"
	rpcCache "dbcache/cluster/network/rpc/types/cache"
	"log"
)

type GetResponse struct {
	Key   string
	Value rpcCache.RpcCacheValue
	Ok    bool
}

type GetRequest struct {
	Key string
}

func (n *RpcNode) Get(key string) (cacher.CacheValue, error) {
	r := new(GetResponse)
	err := n.client.Call(n.rpcAction("RpcGet"), GetRequest{Key: key}, &r)
	return r.Value.Value, err
}

func (n *RpcNode) RpcGet(r GetRequest, resp *GetResponse) error {
	resp.Key = r.Key
	var val cacher.CacheValue
	val, resp.Ok = hostNetwork.cache.Get(r.Key)
	resp.Value = rpcCache.RpcCacheValue{Value: val}
	return nil
}

type SetResponse struct {
}

type SetRequest struct {
	Key   string
	Value rpcCache.RpcCacheValue
}

func (n *RpcNode) Set(key string, value cacher.CacheValue) error {
	r := new(SetResponse)
	rv := rpcCache.RpcCacheValue{Value: value}
	err := n.client.Call(n.rpcAction("RpcSet"), SetRequest{Key: key, Value: rv}, &r)
	return err
}

func (n *RpcNode) RpcSet(r SetRequest, resp *SetResponse) error {
	err := hostNetwork.cache.Set(r.Key, r.Value.Value)
	log.Println("setting value ", r.Key, r.Value.Value)
	return err
}
