package committer

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
)

type VersionCompareComitter struct{}

func (vc *VersionCompareComitter) Commit(c cacher.Cache, b buffer.Buffer) error {
	if b.IsEmpty() {
		return nil
	}
	for k, v := range b.All() {
		val, ok := c.Get(k)
		if !ok || v.Version() > val.Version() {
			err := c.Set(k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateVersionCompareCommitter() *VersionCompareComitter {
	return &VersionCompareComitter{}
}
