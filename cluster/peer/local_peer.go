package peer

import (
	"dbcache/cluster/partition"
)

type LocalPeer struct {
	name      string
	port      uint16
	partition partition.Partition
}

func CreateLocalPeer(name string, port uint16, part partition.Partition) Peer {
	return LocalPeer{name: name, port: port, partition: part}
}

func (p LocalPeer) Name() string {
	return p.name
}

func (p LocalPeer) GePort() uint16 {
	return p.port
}

func (p LocalPeer) Port() uint16 {
	return p.port
}

func (p LocalPeer) GetPartition() partition.Partition {
	return p.partition
}

func (p LocalPeer) Partition() partition.Partition {
	return p.partition
}

func (p LocalPeer) IsSameAs(other Peer) bool {
	return p.Name() == other.Name() && p.Port() == other.Port()
}

func (p LocalPeer) SetPartition(part partition.Partition) Peer {
	p.partition = part
	return p
}

func (p LocalPeer) SetPort(port uint16) Peer {
	p.port = port
	return p
}
