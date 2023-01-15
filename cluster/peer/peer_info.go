package peer

import "dbcache/cluster/version"

type PeerInfo interface {
	Version() version.Version
	IsAlive() bool
	MarkAsDead() PeerInfo
}
