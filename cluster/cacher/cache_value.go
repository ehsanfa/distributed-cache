package cacher

type CacheVersion int

type CacheValue struct {
	Value   string
	Version CacheVersion
}

func NewCacheValue() CacheValue {
	return CacheValue{}
}

func (v CacheVersion) touch() CacheVersion {
	v += 1
	return v
}
