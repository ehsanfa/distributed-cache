package cluster

import (
	"dbcache/cluster/cacher"
)

type CacheRequest struct {
	Action  int8
	Key     string
	Value   string
	Version cacher.CacheVersion
}

type GetCacheResponse struct {
	Ok    bool
	Value string
}

type PutCacheResponse bool

type ShareCacheResposne struct {
	Cache map[string]cacher.CacheValue
}

func (n *Node) Get(key string, resp *GetCacheResponse) error {
	val, ok := thisNode.cache.Get(key)
	*resp = GetCacheResponse{ok, val.Value}
	return nil
}

func (n *Node) Put(req CacheRequest, resp *PutCacheResponse) error {
	v, _ := thisNode.cache.Get(req.Key)
	thisNode.bufferize(cacher.CacheEntity{
		req.Key,
		cacher.CacheValue{req.Value, v.Version},
	})
	*resp = true
	return nil
}

func (n *Node) ShareCache(req cacher.CacheEntity, resp *ShareCacheResposne) error {
	*resp = ShareCacheResposne{Cache: thisNode.cache.All()}
	return nil
}
