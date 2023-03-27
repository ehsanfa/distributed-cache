package rpc

func (n *RpcNode) Ping() (bool, error) {
	var resp PingResponse
	err := n.client.Call("RpcNetwork.RpcPing", PingRequest{}, &resp)
	if err != nil {
		return false, err
	}
	return true, nil
}

type PingRequest struct{}

type PingResponse struct{}

func (n *RpcNetwork) RpcPing(req PingRequest, resp *PingResponse) error {
	return nil
}
