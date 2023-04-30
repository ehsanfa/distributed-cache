package rpc

import (
	"bytes"
	rpcPeer "dbcache/cluster/network/rpc/types/peer"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"encoding/gob"
	"log"
)

type ClusterInfoResponse struct {
	ClusterInfo map[peer.Peer]peer.PeerInfo
}

func (n *RpcNode) GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error) {
	resp := new(ClusterInfoResponse)
	err := n.client.Call(n.rpcAction("RpcGetClusterInfo"), struct{}{}, &resp)
	log.Println("clusterinfo", resp.ClusterInfo)
	return resp.ClusterInfo, err
}

func (n *RpcNode) RpcGetClusterInfo(p struct{}, resp *ClusterInfoResponse) error {
	log.Println("responding to getClusterInfo", hostNetwork.info.All())
	*resp = ClusterInfoResponse{hostNetwork.info.All()}
	return nil
}

type marshalClusterInfoResponse struct {
	Response []marshalClusterInfo
}

type marshalClusterInfo struct {
	Peer     []byte
	PeerInfo []byte
}

func (c *ClusterInfoResponse) MarshalBinary() (data []byte, err error) {
	var mcir []marshalClusterInfo
	for p, pi := range c.ClusterInfo {
		rpcp := rpcPeer.Peer{Peer: p}
		mp, err := rpcp.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		rpcpi := rpcPeer.PeerInfo{Pi: pi}
		mpi, err := rpcpi.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		mci := marshalClusterInfo{mp, mpi}
		mcir = append(mcir, mci)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(marshalClusterInfoResponse{
		Response: mcir,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *ClusterInfoResponse) UnmarshalBinary(data []byte) error {
	r := make(map[peer.Peer]peer.PeerInfo)
	mcir := &marshalClusterInfoResponse{}
	reader := bytes.NewReader(data)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(&mcir); err != nil {
		return err
	}
	ps := rpcPeer.Peer{Peer: peer.CreateLocalPeer("", 0)}
	vers := version.CreateGenClockVersion(8000)
	pis := peer.CreateSimplePeerInfo(vers, true)
	for _, v := range mcir.Response {
		if e := ps.UnmarshalBinary(v.Peer); e != nil {
			return e
		}
		rpcpi := rpcPeer.PeerInfo{Pi: pis}
		if e := rpcpi.UnmarshalBinary(v.PeerInfo); e != nil {
			return e
		}
		a := peer.CreateLocalPeer(ps.Peer.Name(), ps.Peer.Port())
		a = a.SetPartition(ps.Peer.Partition())
		b := peer.CreateSimplePeerInfo(rpcpi.Pi.Version(), rpcpi.Pi.IsAlive())
		r[a] = b
	}
	c.ClusterInfo = r
	return nil
}
