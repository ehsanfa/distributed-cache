package request_handler

import (
	partition "dbcache/partitioning"
	"fmt"
	"net/rpc"
	"time"
)

func (c *Cluster) getNodes(p partition.Partition) *ClusterNodes {
	part := c.getNearestPartition(p)
	peers, ok := c.nodes[part]
	if !ok {
		return NewClusterNodes()
	}
	return peers
}

func (c *Cluster) pickNode(key string) *Peer {
	part := partition.GetPartition(key)
	nodes := c.getNodes(part)
	peer, err := nodes.dequeue()
	if err == nil {
		return peer
	}
	return &c.seeder
}

func NewCluster() *Cluster {
	c := new(Cluster)
	c.nodes = make(map[partition.Partition]*ClusterNodes)

	return c
}

func (c *Cluster) addPeer(part partition.Partition, peer *Peer) {
	if c.nodes == nil {
		c.nodes = make(map[partition.Partition]*ClusterNodes)
	}
	if c.nodes[part] == nil {
		c.nodes[part] = NewClusterNodes()
	}
	if _, ok := c.nodes[part]; !ok {
		c.nodes[part] = NewClusterNodes()
	}
	if !c.nodes[part].exists(peer.info) {
		err := peer.prepare()
		if err != nil {
			fmt.Println(err, peer.info)
			return
		}
		c.nodes[part].add(peer.info, peer)
		fmt.Println("got commit adding part node", part, peer)
		go peer.listen()
	}
}

func (c *Cluster) connect(peer Peer) (*rpc.Client, error) {
	conn, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", peer.info.Name, peer.info.Port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *Cluster) getInfo(infoReceived chan<- bool) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		p := c.seeder
		select {
		case <-ticker.C:
			var resp ShareInfoResponse
			req := ShareCacheRequest{}
			conn, err := c.connect(p)
			defer conn.Close()
			if err != nil {
				fmt.Println(err)
				infoReceived <- false
				continue
			}
			conn.Call("Node.ShareInfo", req, &resp)
			c.info = resp.Info
			for part, peers := range resp.Partitions {
				for pi, _ := range peers {
					mn := &Peer{info: PeerInfo{Name: pi.Name, Port: pi.Port}}
					c.addPeer(part, mn)
				}
			}
			c.sortPartitions()
			infoReceived <- true
		}
	}
}
