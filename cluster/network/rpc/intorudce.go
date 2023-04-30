package rpc

import (
	"dbcache/cluster/peer"
	"dbcache/cluster/version"
	"log"
)

type IntroductionResponse struct{}

type ReqPeer struct {
	Name string
	Port uint16
}

func (n *RpcNode) Introduce(p peer.Peer) error {
	resp := new(IntroductionResponse)
	reqP := ReqPeer{Name: p.Name(), Port: p.Port()}
	err := n.client.Call(n.rpcAction("RpcIntroduce"), reqP, &resp)
	return err
}

func (n *RpcNode) RpcIntroduce(req ReqPeer, r *IntroductionResponse) error {
	ver := version.CreateGenClockVersion(0)
	p := peer.CreateLocalPeer(req.Name, req.Port)
	hostNetwork.info.Add(p, peer.CreateSimplePeerInfo(ver, true))
	log.Println("got intoruction request", p.Name(), p.Port())
	return nil
}
