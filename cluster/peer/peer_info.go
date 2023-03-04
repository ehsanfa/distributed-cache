package peer

import (
	"dbcache/cluster/version"
	"encoding"
)

type PeerInfo interface {
	Version() version.Version
	IsAlive() bool
	MarkAsDead() PeerInfo
	encoding.BinaryMarshaler
}
