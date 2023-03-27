package cacher

type Cachable interface {
	GetKey() string
	GetValue() CacheValue
}

// type CacheEntity struct {
// 	Key   string
// 	Value CacheValue
// }

// func (c *CacheEntity) GetKey() string {
// 	return c.Key
// }

// func (c *CacheEntity) GetValue() CacheValue {
// 	return c.Value
// }
