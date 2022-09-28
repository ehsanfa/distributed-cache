package cluster

import (
	"fmt"
	"time"
	"dbcache/types"
)

type CacheRequest struct {
	Action int8
	Key string
	Value string
	AlreadyAware  map[Peer]bool
}

type CacheResponse types.Resp

type ShareCacheResposne struct {
	Cache    map[string]string
}

type ShareCacheRequest struct {}

func (n *Node) shareWithPeers(req CacheRequest) {
	peers := n.getPeersForPartitioning(req.Key)
	if len(peers) == 0 {
		peers = n.convertInfoToPartitionPeers()
	}
	for peer, _ := range peers {
		if _, ok := req.AlreadyAware[peer]; ok || !peer.isAlive() {
			continue
		}
		go n.share(peer, req)
	}
}

func (n *Node) share(p Peer, req CacheRequest) {
	c, err := n.getConnection(p)
	if err != nil {
		fmt.Println("peer not responding. unbuddying", p, n.peersToGossip())
		n.unbuddy(p)
		return
	}
	defer c.Close()

	fmt.Println("sharing cache with ", p, req, n.getPeersForPartitioning(req.Key))

	var resp CacheResponse
	timer := time.NewTimer(gossipTimeout)
	call := c.Go("Node.Put", req, &resp, nil)
	select{
	case <-call.Done:
		timer.Stop()
	case <-timer.C:
		fmt.Println("peer not responding. leaving it for the gossip to handle TWO", p, n.peersToGossip())
		time.Sleep(2*time.Second)
		return
	}
}

func (n *Node) Get(req CacheRequest, resp *CacheResponse) error {
	thisNode.mu.RLock()
	val, ok := thisNode.cache[req.Key]
	thisNode.mu.RUnlock()
	*resp = CacheResponse{ok, req.Key, val}
	return nil
}

func (n *Node) ShareCache(req ShareCacheRequest, resp *ShareCacheResposne) error {
	*resp = ShareCacheResposne{Cache: thisNode.cache}
	return nil
}

func (n *Node) Put(req CacheRequest, resp *CacheResponse) error {
	// fmt.Println("setting cache for ", req)
	thisNode.mu.Lock()
	thisNode.cache[req.Key] = req.Value
	thisNode.mu.Unlock()
	if req.AlreadyAware == nil {
		req.AlreadyAware = make(map[Peer]bool)
	}
	req.AlreadyAware[*thisNode.getPeer()] = true
	// thisNode.shareWithPeers(req)
	*resp = CacheResponse{true, req.Key, ""}
	return nil
}