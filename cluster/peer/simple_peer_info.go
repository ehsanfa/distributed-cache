package peer

import "dbcache/cluster/version"

type SimplePeerInfo struct {
	version version.Version
	isAlive bool
}

func CreateSimplePeerInfo(ver version.Version, isAlive bool) SimplePeerInfo {
	return SimplePeerInfo{version: ver, isAlive: isAlive}
}

func (i SimplePeerInfo) MarkAsDead() PeerInfo {
	i.isAlive = false
	return i
}

func (i SimplePeerInfo) Version() version.Version {
	return i.version
}

func (i SimplePeerInfo) IsAlive() bool {
	return i.isAlive
}

func (i SimplePeerInfo) MarshalBinary() (data []byte, err error) {
	return []byte("sda"), nil
}
