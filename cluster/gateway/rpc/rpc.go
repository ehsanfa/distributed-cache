package rpc

import (
	"dbcache/cluster/info"
	"dbcache/cluster/network"
)

type RpcGateway struct {
	clusterInfo info.ClusterInfo
	network     network.Network
}

func CreateRpcGateway(ci info.ClusterInfo, n network.Network) *RpcGateway {
	return &RpcGateway{clusterInfo: ci, network: n}
}
