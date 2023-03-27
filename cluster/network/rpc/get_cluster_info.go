package rpc

import (
	"bytes"
	rpcPeer "dbcache/cluster/network/rpc/types/peer"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"encoding/gob"
)

type ClusterInfoResponse struct {
	ClusterInfo map[peer.Peer]peer.PeerInfo
}

func (n *RpcNode) GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error) {
	resp := new(ClusterInfoResponse)
	err := n.client.Call("RpcNode.RpcGetClusterInfo", struct{}{}, &resp)
	return resp.ClusterInfo, err
}

func (n *RpcNode) RpcGetClusterInfo(p struct{}, resp *ClusterInfoResponse) error {
	*resp = ClusterInfoResponse{hostNetwork.info.GetClusterInfo()}
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
	parts := partition.CreateSimplePartition("")
	ps := rpcPeer.Peer{Peer: peer.CreateLocalPeer("", 0, &parts)}
	vers := version.CreateGenClockVersion(0)
	pis := peer.CreateSimplePeerInfo(vers, true)
	for _, v := range mcir.Response {
		if e := ps.UnmarshalBinary(v.Peer); e != nil {
			return e
		}
		rpcpi := rpcPeer.PeerInfo{Pi: pis}
		if e := rpcpi.UnmarshalBinary(v.PeerInfo); e != nil {
			return e
		}
		a := peer.CreateLocalPeer(ps.Peer.Name(), ps.Peer.Port(), ps.Peer.Partition())
		b := peer.CreateSimplePeerInfo(pis.Version(), pis.IsAlive())
		r[a] = b
	}
	c.ClusterInfo = r
	return nil
}
