package cluster

import (
	// "fmt"
	// "time"
)

// var counter int64 = 0

// func (n *Node) reportCount() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			fmt.Println("report counter", counter/5, n.partition)
// 			counter = 0
// 		}
// 	}
// }

type CacheRequest struct {
	Action int8
	Key string
	Value string
	Version CacheVersion
}

type GetCacheResponse struct {
	Ok bool
	Value string
}

type PutCacheResponse bool

type ShareCacheResposne struct {
	Cache    map[string]CacheValue
}

type CacheEntity struct {
	Key        string
	Value      CacheValue
}

func (n *Node) Get(key string, resp *GetCacheResponse) error {
	thisNode.cacheMu.RLock()
	val, ok := thisNode.cache[key]
	thisNode.cacheMu.RUnlock()
	*resp = GetCacheResponse{ok, val.Value}
	return nil
}

func NewCacheValue() CacheValue {
	return CacheValue{}
}

// func (c CacheValue) update(val CacheValue) CacheValue {
// 	c.Value = val
// 	c.Version = c.Version.update()
// 	return c
// }

func (v CacheVersion) update() CacheVersion {
	v += 1
	return v
}

func (n *Node) getCacheVersion(key string) CacheVersion {
	thisNode.cacheVersionsMu.RLock()
	ver, ok := thisNode.cacheVersions[key]
	thisNode.cacheVersionsMu.RUnlock()
	if !ok {
		return 0
	}
	return ver
}

func (n *Node) updateVersion(key string, version CacheVersion) {
	thisNode.cacheVersionsMu.Lock()
	thisNode.cacheVersions[key] = version
	thisNode.cacheVersionsMu.Unlock()
}

func (n *Node) put(key string, value CacheValue) {
	thisNode.cacheMu.Lock()
	if _, ok := thisNode.cache[key]; !ok {
		thisNode.cache[key] = NewCacheValue()
	}
	thisNode.cache[key] = value
	thisNode.cacheMu.Unlock()
	n.updateVersion(key, value.Version)
}

func (n *Node) Put(req CacheRequest, resp *PutCacheResponse) error {
	v := thisNode.getCacheVersion(req.Key)
	// go thisNode.put(req.Key, req.Value)
	thisNode.bufferize(CacheEntity{
		req.Key,
		CacheValue{req.Value,v},
	})
	*resp = true
	return nil
}

func (n *Node) ShareCache(req CacheEntity, resp *ShareCacheResposne) error {
	*resp = ShareCacheResposne{Cache: thisNode.cache}
	return nil
}