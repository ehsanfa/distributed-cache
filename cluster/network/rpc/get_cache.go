package network

import "dbcache/cluster/cacher"

func (n *RpcNode) GetCache() (map[string]cacher.CacheValue, error) {
	return nil, nil
}
