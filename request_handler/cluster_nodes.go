package request_handler

import (
	"dbcache/request_handler/deque"
)

type ClusterNodes struct {
	nodes  map[PeerInfo]*Peer 
	deq *deque.Deque
}

type hasPeer interface {
	getPeer() *Peer
}

func (p *Peer) getPeer() *Peer {
	return p
}

func (c *ClusterNodes) exists(pi PeerInfo) bool {
	_, ok := c.nodes[pi]
	return ok
}

func (c *ClusterNodes) add(pi PeerInfo, p *Peer) {
	c.nodes[pi] = p
	c.deq.Enqueue(p)
}

func (c *ClusterNodes) dequeue() (*Peer, error) {
	p, err := c.deq.Dequeue()
	if err != nil {
		return nil, err
	}
	return p.(hasPeer).getPeer(), nil
}

func (c *ClusterNodes) isEmpty() bool {
	return c.deq.Count() == 0
}

func NewClusterNodes() *ClusterNodes {
	c := new(ClusterNodes)
	c.nodes = make(map[PeerInfo]*Peer)
	c.deq = deque.NewDeque()
	return c
}