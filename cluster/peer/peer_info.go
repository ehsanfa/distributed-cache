package peer

import (
	"dbcache/cluster/version"
)

type PeerInfo interface {
	Version() version.Version
	IsAlive() bool
	MarkAsDead() PeerInfo
	MarkAsAlive() PeerInfo
	Type() PeerType
	IsCacher() bool
}

type PeerType uint8

const (
	Unspecified PeerType = iota
	Seeder
	Cacher
	Gateway
)
