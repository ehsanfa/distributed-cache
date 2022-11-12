package cluster

import (
	"fmt"
	"time"
	"math/rand"
	// "dbcache/types"
	ll "github.com/ehsanfa/linked-list"
)

var counter int64 = 0
const shareMin = 2000
const shareMax = 5000
var shareBufferInterval = getShareBufferInterval()
const bufferSizeLimit = 1 << 32

func getShareBufferInterval() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(shareMax - shareMin) + shareMin) * time.Millisecond
}

func (n *Node) reportCount() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("report counter", counter/5, n.partition)
			counter = 0
		}
	}
}

type CacheRequest struct {
	Action int8
	Key string
	Value string
	Version CacheVersion
}

type CacheResponse struct {
	Ok bool
	Key string
	Value string
}

type ShareCacheResposne struct {
	Cache    map[string]CacheValue
}

type CacheEntity struct {
	Key        string
	Value      CacheValue
}

type Buffer struct {
	internal      ll.LinkedList
	SharingBuffer map[string]CacheValue
}

func (b *Buffer) isEmpty() bool {
	return b.internal.Count() == 0
}

func (b *Buffer) shareIsEmpty() bool {
	return len(b.SharingBuffer) == 0
}

func (b *Buffer) add(c CacheEntity) {
	b.internal.Append(c)
	if b.SharingBuffer == nil {
		b.SharingBuffer = make(map[string]CacheValue)
	}
	b.SharingBuffer[c.Key] = c.Value
}

func (b *Buffer) count() int {
	return b.internal.Count()
}

type ShareBufferRequest struct {
	Buffer        Buffer
	AlreadyAware  map[Peer]bool
}

func (n *Node) startCleaningBuffer() {
	ticker := time.NewTicker(shareBufferInterval)
	fmt.Println("cleaning buffer every", shareBufferInterval)
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
	if n.buffer.shareIsEmpty() {
		return
	}
	peers := n.getPeersToShareBuffer()
	// fmt.Println("sharing buffer with peers", peers)
	req := ShareBufferRequest{Buffer: n.buffer}
	req.AlreadyAware = make(map[Peer]bool)
	req.AlreadyAware[*n.getPeer()] = true
	for peer, _ := range peers {
		go n.share(peer, req)
	}
	n.resetBuffer()
}

func (n *Node) resetBuffer() {
	n.buffer = Buffer{}
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
	if p == *n.getPeer() {
		return
	}
	fmt.Println("sharing buffer with peer", p, req)
	c, err := n.getConnection(p)
	if err != nil {
		fmt.Println("peer not responding. unbuddying", p, n.peersToGossip())
		n.unbuddy(p)
		return
	}

	var resp CacheResponse
	timer := time.NewTimer(gossipTimeout)
	call := c.Go("Node.ShareBuffer", req, &resp, nil)
	if call.Error != nil {
		fmt.Println(call.Error)
	}
	select{
	case <-call.Done:
		timer.Stop()
	case <-timer.C:
		fmt.Println("peer not responding. leaving it for the gossip to handle TWO", p, n.peersToGossip())
		time.Sleep(2*time.Second)
		return
	}
}

func (n *Node) Get(key string, resp *CacheResponse) error {
	counter++
	thisNode.cacheMu.RLock()
	val, ok := thisNode.cache[key]
	thisNode.cacheMu.RUnlock()
	*resp = CacheResponse{ok, key, val.Value}
	return nil
}

func (n *Node) ShareBuffer(req ShareBufferRequest, resp *CacheResponse) error {
	fmt.Println("heeeeeeeeeeeeere")
	*resp = CacheResponse{}
	req.AlreadyAware[*thisNode.getPeer()] = true
	go thisNode.handOverBuffer(req)
	for k, v := range req.Buffer.SharingBuffer {
		thisNode.buffer.internal.Append(CacheEntity{k, v})
	}
	return nil
}

func (n *Node) commit() {
	if n.buffer.isEmpty() {
		return
	}
	buffer := n.buffer.internal
	for buffer.Count() != 0 {
		if buffer.Tail() == nil {
			return
		}
		c := buffer.Pop()
		e := c.Value().(CacheEntity)
		if e.Value.Version > n.getCacheVersion(e.Key) {
			n.put(e.Key, CacheValue{e.Value.Value, e.Value.Version})
		}
	}
}

func NewCacheValue() CacheValue {
	return CacheValue{}
}

// func (c CacheValue) update(val CacheValue) CacheValue {
// 	c.Value = val
// 	c.Version = c.Version.update()
// 	return c
// }

func (v CacheVersion) update() CacheVersion {
	v += 1
	return v
}

func (n *Node) getCacheVersion(key string) CacheVersion {
	thisNode.cacheMu.RLock()
	val, ok := thisNode.cache[key]
	thisNode.cacheMu.RUnlock()
	if !ok {
		return 0
	}
	return val.Version
}

func (n *Node) put(key string, value CacheValue) {
	thisNode.cacheMu.Lock()
	if _, ok := thisNode.cache[key]; !ok {
		thisNode.cache[key] = NewCacheValue()
	}
	thisNode.cache[key] = value
	thisNode.cacheMu.Unlock()
	fmt.Println(thisNode.cache[key])
}

func (n *Node) addToBuffer(req CacheEntity) {
	// TODO; use queue instead of array
	req.Value.Version = req.Value.Version.update()
	n.buffer.add(req)
	if n.buffer.count() > bufferSizeLimit {
		n.bufferSizeExceeded <- true
	}
	fmt.Println("added to buffer", req)
	// fmt.Println("added to buffer", n.buffer.Buffer.Count())
}

func (n *Node) Put(req CacheRequest, resp *CacheResponse) error {
	counter++
	v := thisNode.getCacheVersion(req.Key)
	// go thisNode.put(req.Key, req.Value)
	thisNode.addToBuffer(CacheEntity{
		req.Key,
		CacheValue{req.Value,v},
	})
	*resp = CacheResponse{true, req.Key, ""}
	return nil
}

func (n *Node) ShareCache(req CacheEntity, resp *ShareCacheResposne) error {
	*resp = ShareCacheResposne{Cache: thisNode.cache}
	return nil
}