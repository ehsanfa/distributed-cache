package rpc

import (
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
)

type IntroductionResponse struct{}

type ReqPeer struct {
	Name string
	Port uint16
	Type peer.PeerType
}

func (n *RpcNode) Introduce(peerType peer.PeerType, p peer.Peer) error {
	resp := new(IntroductionResponse)
	reqP := ReqPeer{Name: p.Name(), Port: p.Port(), Type: peerType}
	err := n.client.Call(n.rpcAction("RpcIntroduce"), reqP, &resp)
	return err
}

func (n *RpcNode) RpcIntroduce(req ReqPeer, r *IntroductionResponse) error {
	p := peer.CreateLocalPeer(req.Name, req.Port)
	peerInfo, ok := hostNetwork.info.Get(p)
	if ok {
		peerInfo = peerInfo.MarkAsAlive()
	} else {
		ver := version.CreateGenClockVersion(0)
		peerInfo = peer.CreateSimplePeerInfo(req.Type, ver, true)
	}
	hostNetwork.info.Add(p, peerInfo)
	return nil
}
