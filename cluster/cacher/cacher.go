package cacher

type Cache interface {
	Get(key string) (value CacheValue, ok bool)
	Exists(key string) bool
	Set(key string, value CacheValue) error
	All() map[string]CacheValue
	Touch(key string)
	Version(key string) CacheVersion
	Replace(map[string]CacheValue)
	Delete(key string)
}

type CachStore interface {
	GetCache() Cache
	SetCache(Cache)
}

// type PutCacheResponse bool

// type ShareCacheResposne struct {
// 	Cache map[string]CacheValue
// }
