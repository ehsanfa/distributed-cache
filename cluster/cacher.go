package cluster

import (
	"fmt"
	"time"
	"dbcache/types"
	ll "github.com/ehsanfa/linked-list"
)

var counter int64 = 0
const shareBufferInterval = 5 * time.Second
const bufferSizeLimit = 1 << 32

func (n *Node) reportCount() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("report counter", counter/5)
			counter = 0
		}
	}
}

type CacheRequest struct {
	Action int8
	Key string
	Value string
}

type CacheResponse types.Resp

type ShareCacheResposne struct {
	Cache    map[string]string
}

type ShareCacheRequest struct {
	Key        string
	Value      string
	UpdatedAt  time.Time
}

type SharingBuffer struct {
	Buffer      ll.LinkedList
}

func (b *SharingBuffer) isEmpty() bool {
	return b.Buffer.Count() == 0
}

type ShareBufferRequest struct {
	Buffer        SharingBuffer
	AlreadyAware  map[Peer]bool
}

func (n *Node) startCleaningBuffer() {
	ticker := time.NewTicker(shareBufferInterval)
	for {
		select {
		case <-ticker.C:
			n.commit()
			n.shareBuffer()
		case <-n.bufferSizeExceeded:
			n.commit()
			n.shareBuffer()
		}
	}
}

func (n *Node) getPeersToShareBuffer() map[Peer]bool {
	peers := n.getPeersForPartitioning()
	if len(peers) == 0 {
		peers = n.convertInfoToPartitionPeers()
	}
	return peers
}

func (n *Node) shareBuffer() {
	if n.shareCacheBuffer.isEmpty() {
		return
	}
	peers := n.getPeersToShareBuffer()
	// fmt.Println("sharing buffer with peers", peers)
	req := ShareBufferRequest{Buffer: n.shareCacheBuffer}
	req.AlreadyAware = make(map[Peer]bool)
	req.AlreadyAware[*n.getPeer()] = true
	for peer, _ := range peers {
		go n.share(peer, req)
	}
	n.resetBuffer()
}

func (n *Node) resetBuffer() {
	n.shareCacheBuffer = SharingBuffer{}
}

func (n *Node) handOverBuffer(req ShareBufferRequest) {
	peers := n.getPeersToShareBuffer()
	for peer, _ := range peers {
		if _, ok := req.AlreadyAware[peer]; ok || !peer.isAlive() {
			continue
		}
		go n.share(peer, req)
	}
}

func (n *Node) share(p Peer, req ShareBufferRequest) {
	c, err := n.getConnection(p)
	if err != nil {
		fmt.Println("peer not responding. unbuddying", p, n.peersToGossip())
		n.unbuddy(p)
		return
	}
	defer c.Close()

	var resp CacheResponse
	timer := time.NewTimer(gossipTimeout)
	call := c.Go("Node.ShareBuffer", req, &resp, nil)
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
	counter++
	thisNode.cacheMu.RLock()
	val, ok := thisNode.cache[req.Key]
	thisNode.cacheMu.RUnlock()
	*resp = CacheResponse{ok, req.Key, val}
	return nil
}

func (n *Node) ShareBuffer(req ShareBufferRequest, resp *ShareCacheResposne) error {
	*resp = ShareCacheResposne{Cache: thisNode.cache}
	req.AlreadyAware[*thisNode.getPeer()] = true
	go thisNode.handOverBuffer(req)
	// for _, c := range req.Buffer.Buffer {
	// 	thisNode.put(c.Key, c.Value)
	// }
	return nil
}

func (n *Node) commit() {
	if n.shareCacheBuffer.isEmpty() {
		return
	}
	for n.shareCacheBuffer.Buffer.Count() != 0 {
		if n.shareCacheBuffer.Buffer.Tail() == nil {
			return
		}
		node := n.shareCacheBuffer.Buffer.Pop()
		val := node.Value().(ShareCacheRequest)
		n.put(val.Key, val.Value)
	}
}

func (n *Node) put(key, value string) {
	thisNode.cacheMu.Lock()
	thisNode.cache[key] = value
	thisNode.cacheMu.Unlock()
}

func (n *Node) addToBuffer(req ShareCacheRequest) {
	// TODO; use queue instead of array
	n.shareCacheBuffer.Buffer.Append(req)
	if n.shareCacheBuffer.Buffer.Count() > bufferSizeLimit {
		n.bufferSizeExceeded <- true
	}
	// fmt.Println("added to buffer", n.shareCacheBuffer.Buffer.Count())
}

func (n *Node) Put(req CacheRequest, resp *CacheResponse) error {
	counter++
	// go thisNode.put(req.Key, req.Value)
	thisNode.addToBuffer(ShareCacheRequest{
		req.Key,
		req.Value,
		time.Now(),
	})
	*resp = CacheResponse{true, req.Key, ""}
	return nil
}