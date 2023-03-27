package cacher

type CacheValue interface {
	GetValue() string
	SetValue(string)
	Version() int
	IncrementVersion()
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}
