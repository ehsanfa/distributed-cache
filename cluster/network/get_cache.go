package network

import (
	"dbcache/cluster/cacher"
)

type CacheProvider interface {
	GetCache() map[string]cacher.CacheValue
}
