package request_handler

import (
	"fmt"
	"time"
	"net/rpc"
	partition "dbcache/partitioning"
)

func (c *Cluster) getNodes(p partition.Partition) *ClusterNodes {
	p = c.getNearestPartition(p)
	peers, ok := c.nodes[p]
	if !ok {
		return NewClusterNodes()
	}
	return peers
}

func (c *Cluster) pickNode(key string) *Peer {
	part := partition.GetPartition(key)
	nodes := c.getNodes(part)
	if !nodes.isEmpty() {
		if peer, err := nodes.dequeue(); err == nil {
			return peer
		}
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
		// c.nodes[part][peer.info] = peer
		fmt.Println("commit adding part node", part, peer)
		go peer.listen()
	}
}

func (c *Cluster) getInfo(infoReceived chan<- bool){
	p := c.seeder
	conn, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.info.Name, p.info.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			// fmt.Println("getting info")
			var resp ShareInfoResponse
			req := ShareCacheRequest{}
			conn.Call("Node.ShareInfo", req, &resp)
			c.info = resp.Info
			// fmt.Println("got info ", resp)
			// var pps map[*Peer]bool
			for part, peers := range resp.Partitions {
				// pps = make(/[*Peer]bool)
				for pi, _ := range peers {
					mn := &Peer{info: PeerInfo{Name: pi.Name, Port: pi.Port}}
					c.addPeer(part, mn)
				}
				// c.nodes[part] = pps
			}
			// c.nodes = resp.Partitions
			c.sortPartitions()
			infoReceived <- true
		}
	}
}