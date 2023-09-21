package rpc

import (
	"dbcache/cluster/buffer"
)

type RpcUpdateBufferResp struct {
	Status bool
}

type UpdateBufferRequest struct {
	buff buffer.Buffer
}

func (n *RpcNode) UpdateBuffer(source buffer.Buffer) error {
	req := new(UpdateBufferRequest)
	req.buff = source
	resp := new(RpcUpdateBufferResp)
	err := n.client.Call(n.rpcAction("RpcUpdateBuffer"), req, &resp)
	return err
}

func (n *RpcNode) RpcUpdateBuffer(req UpdateBufferRequest, resp *RpcUpdateBufferResp) error {
	resp.Status = true
	hostNetwork.buffer.Merge(req.buff)
	return nil
}
