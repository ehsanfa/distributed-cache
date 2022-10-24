package request_handler

import (
	"fmt"
	"time"
	"net/rpc"
	partition "dbcache/partitioning"
)

func (c *Cluster) getNodes(p partition.Partition) map[*Peer]*rpc.Client {
	p = c.getNearestPartition(p)
	peers, ok := c.partitions[p]
	if !ok {
		return make(map[*Peer]*rpc.Client)
	}
	return peers
}

func (c *Cluster) pickNode(key string) *Peer {
	part := partition.GetPartition(key)
	for peer, _ := range c.getNodes(part) {
		return peer
	}
	return &c.seeder
}

func NewCluster() *Cluster {
	c := new(Cluster)
	c.partitions = make(map[partition.Partition]map[*Peer]*rpc.Client)
	return c
}

func (c *Cluster) addPeer(part partition.Partition, peer *Peer) {
	if c.partitions == nil {
		c.partitions = make(map[partition.Partition]map[*Peer]*rpc.Client)
	}
	if c.partitions[part] == nil {
		c.partitions[part] = make(map[*Peer]*rpc.Client)
	}
	if _, ok := c.partitions[part]; !ok {
		c.partitions[part] = make(map[*Peer]*rpc.Client)
	}
	if _, ok := c.partitions[part][peer]; !ok {
		conn, err := dial(peer)
		if err != nil {
			return
		}
		peer.conn = conn
		c.partitions[part][peer] = conn
		go peer.listen()
	}
}

func (c *Cluster) getInfo(infoReceived chan<- bool){
	p := c.seeder
	conn, err := rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
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
				for pe, _ := range peers {
					mn := &Peer{Name: pe.Name, Port: pe.Port}
					c.addPeer(part, mn)
				}
				// c.partitions[part] = pps
			}
			// c.partitions = resp.Partitions
			c.sortPartitions()
			infoReceived <- true
		}
	}
}