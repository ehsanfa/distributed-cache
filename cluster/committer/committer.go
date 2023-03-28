package committer

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
)

type Committer interface {
	Commit(cacher.Cache, buffer.Buffer) error
}
