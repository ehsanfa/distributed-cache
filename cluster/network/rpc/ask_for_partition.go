package rpc

import "dbcache/cluster/partition"

// Not necessary for now. Partitions can be manually assigned
func (n *RpcNode) AskForParition() (partition.Partition, error) {
	return nil, nil
}
