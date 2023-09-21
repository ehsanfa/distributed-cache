package cacher

type CacheValue interface {
	GetValue() string
	SetValue(string) CacheValue
	Version() int
	IncrementVersion() CacheValue
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}
