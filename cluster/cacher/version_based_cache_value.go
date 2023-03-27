package cacher

type VersionBasedCacheValue struct {
	value   string
	version int
}

func (v *VersionBasedCacheValue) GetValue() string {
	return v.value
}

func (v *VersionBasedCacheValue) SetValue(value string) {
	v.value = value
}

func (v *VersionBasedCacheValue) Version() int {
	return v.version
}

func NewVersionBasedCacheValue(value string, version int) *VersionBasedCacheValue {
	return &VersionBasedCacheValue{value: value, version: version}
}

func (v *VersionBasedCacheValue) IncrementVersion() {
	v.version += 1
}
