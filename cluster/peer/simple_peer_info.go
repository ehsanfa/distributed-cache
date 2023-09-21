package peer

import (
	"dbcache/cluster/version"
)

type SimplePeerInfo struct {
	version  version.Version
	isAlive  bool
	peerType PeerType
}

func CreateSimplePeerInfo(peerType PeerType, ver version.Version, isAlive bool) SimplePeerInfo {
	return SimplePeerInfo{peerType: peerType, version: ver, isAlive: isAlive}
}

func (i SimplePeerInfo) MarkAsDead() PeerInfo {
	i.isAlive = false
	i.version = i.version.Increment()
	return i
}

func (i SimplePeerInfo) MarkAsAlive() PeerInfo {
	i.isAlive = true
	i.version = i.version.Increment()
	return i
}

func (i SimplePeerInfo) Version() version.Version {
	return i.version
}

func (i SimplePeerInfo) IsAlive() bool {
	return i.isAlive
}

func (i SimplePeerInfo) Type() PeerType {
	return i.peerType
}

func (i SimplePeerInfo) IsCacher() bool {
	return i.Type() == Cacher
}
