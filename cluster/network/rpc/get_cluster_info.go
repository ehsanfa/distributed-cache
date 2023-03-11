package network

import (
	"bytes"
	"dbcache/cluster/partition"
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"encoding/gob"
)

func (n *RpcNode) GetClusterInfo() (map[peer.Peer]peer.PeerInfo, error) {
	resp := new(ClusterInfoResponse)
	n.client.Call("RpcNetwork.RpcGetClusterInfo", struct{}{}, &resp)
	return resp.ClusterInfo, nil
}

type ClusterInfoResponse struct {
	ClusterInfo map[peer.Peer]peer.PeerInfo
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
		mp, err := p.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		mpi, err := pi.MarshalBinary()
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
	ps := peer.CreateLocalPeer("hasan", 0, &parts)
	vers := version.CreateGenClockVersion()
	pis := peer.CreateSimplePeerInfo(vers, true)
	for _, v := range mcir.Response {
		if e := ps.UnmarshalBinary(v.Peer); e != nil {
			return e
		}
		if e := pis.UnmarshalBinary(v.PeerInfo); e != nil {
			return e
		}
		a := peer.CreateLocalPeer(ps.Name(), ps.Port(), ps.Partition())
		b := peer.CreateSimplePeerInfo(pis.Version(), pis.IsAlive())
		r[a] = b
	}
	c.ClusterInfo = r
	return nil
}

func (n *RpcNetwork) RpcGetClusterInfo(p struct{}, resp *ClusterInfoResponse) error {
	*resp = ClusterInfoResponse{n.server.info.GetClusterInfo()}
	return nil
}
