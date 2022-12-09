package cluster

import (
	"fmt"
	"math/rand"
	"time"

	ll "github.com/ehsanfa/linked-list"
)

const shareMin = 2000
const shareMax = 5000

var shareBufferInterval = getShareBufferInterval()

const bufferSizeLimit = 1 << 32

func getShareBufferInterval() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(shareMax-shareMin)+shareMin) * time.Millisecond
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

func (n *Node) addToBuffer(c CacheEntity) {
	n.bufferMu.Lock()
	n.buffer.add(c)
	n.bufferMu.Unlock()
}

func (n *Node) isBufferEmpty() bool {
	n.bufferMu.RLock()
	empty := n.buffer.isEmpty()
	n.bufferMu.RUnlock()
	return empty
}

func (n *Node) isSharingBufferEmpty() bool {
	n.bufferMu.RLock()
	empty := n.buffer.shareIsEmpty()
	n.bufferMu.RUnlock()
	return empty
}

func (b *Buffer) count() int {
	return b.internal.Count()
}

func (n *Node) bufferSize() int {
	n.bufferMu.RLock()
	size := n.buffer.count()
	n.bufferMu.RUnlock()
	return size
}

func (n *Node) getBuffer() Buffer {
	n.bufferMu.RLock()
	b := n.buffer
	n.bufferMu.RUnlock()
	return b
}

type ShareBufferRequest struct {
	Buffer       Buffer
	AlreadyAware map[Peer]bool
}

func (n *Node) startCleaningBuffer() {
	ticker := time.NewTicker(shareBufferInterval)
	fmt.Println("cleaning buffer every", shareBufferInterval)
	for {
		select {
		case <-ticker.C:
			n.commit()
			// n.shareBuffer()
		case <-n.bufferSizeExceeded:
			n.commit()
			// n.shareBuffer()
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
	if n.isSharingBufferEmpty() {
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
	n.bufferMu.Lock()
	n.buffer = Buffer{}
	n.bufferMu.Unlock()
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
	c, err := n.getConnection(p)
	if err != nil {
		fmt.Println("peer not responding. unbuddying", p, n.peersToGossip())
		n.unbuddy(p)
		return
	}

	var resp interface{}
	timer := time.NewTimer(gossipTimeout)
	call := c.Go("Node.ShareBuffer", req, &resp, nil)
	if call.Error != nil {
		fmt.Println(call.Error, "here")
	}
	select {
	case <-call.Done:
		timer.Stop()
	case <-timer.C:
		fmt.Println("peer not responding. leaving it for the gossip to handle TWO", p, n.peersToGossip())
		time.Sleep(2 * time.Second)
		return
	}
}

func (n *Node) ShareBuffer(req ShareBufferRequest, resp *interface{}) error {
	*resp = nil
	req.AlreadyAware[*thisNode.getPeer()] = true
	go thisNode.handOverBuffer(req)
	for k, v := range req.Buffer.SharingBuffer {
		thisNode.addToBuffer(CacheEntity{k, v})
	}
	return nil
}

func (n *Node) commit() {
	if n.isBufferEmpty() {
		return
	}
	buffer := n.getBuffer()
	for buffer.internal.Count() != 0 {
		if buffer.internal.Tail() == nil {
			return
		}
		c := buffer.internal.Pop()
		e := c.Value().(CacheEntity)
		if e.Value.Version > n.getCacheVersion(e.Key) {
			n.put(e.Key, CacheValue{e.Value.Value, e.Value.Version})
		}
	}
}

func (n *Node) bufferize(req CacheEntity) {
	req.Value.Version = req.Value.Version.update()
	n.addToBuffer(req)
	if n.bufferSize() > bufferSizeLimit {
		n.bufferSizeExceeded <- true
	}
	// fmt.Println("added to buffer", n.buffer.Buffer.Count())
}
