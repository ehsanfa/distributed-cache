package connection

import (
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"testing"
)

type MockConnection struct {}

func (c MockConnection) Introduce(peer.Peer) (IntroductionResponse, error) {

}

type MockIntroductionResponse struct {
}

func (m MockIntroductionResponse) ClusterInfo() info.ClusterInfo {

}

func (m MockIntroductionResponse) Cache() cacher.Cache {

}

func (m MockIntroductionResponse) Partition() partition.Partition {

}

func TestIntroduce(t *testing.T) {
	part := partition.CreateSimplePartition("0")
	seeder := peer.CreateLocalPeer("seeder", 1234, part)
	p := peer.CreateLocalPeer("peer", 1234, part)
	c := 
}
