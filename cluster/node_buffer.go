package cluster

import (
	"dbcache/cluster/buffer"
	"fmt"
	"math/rand"
	"time"
)

const shareMin = 2000
const shareMax = 5000

const bufferSizeLimit = 1 << 32

func getShareBufferInterval() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(shareMax-shareMin)+shareMin) * time.Millisecond
}

func (n *Node) startCleaningBuffer() {
	ticker := time.NewTicker(getShareBufferInterval())
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

type ShareBufferRequest struct {
	Buffer       buffer.Buffer
	AlreadyAware map[Peer]bool
}

func (n *Node) shareBuffer() {
	if n.buffer.IsEmpty() {
		return
	}
	peers := n.getPeersToShareBuffer()
	if len(peers) == 0 {
		return
	}
	// fmt.Println("sharing buffer with peers", peers)
	req := ShareBufferRequest{Buffer: n.buffer}
	req.AlreadyAware = make(map[Peer]bool)
	req.AlreadyAware[*n.getPeer()] = true
	for peer, _ := range peers {
		if peer.isSame(*n.getPeer()) {
			continue
		}
		go n.share(peer, req)
	}
	n.resetSharingBuffer()
}

func (n *Node) resetSharingBuffer() {
	n.bufferMu.Lock()
	n.buffer = *n.buffer.resetSharingBuffer()
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
	buffer := n.getBuffer()
	if n.isBufferEmpty() {
		return
	}
	for buffer.internal.Count() != 0 {
		if buffer.internal.Tail() == nil {
			return
		}
		c := buffer.internal.Pop()
		e := c.Value().(CacheEntity)
		if e.Value.Version > n.cache.Version(e.Key) {
			n.cache.Set(e.Key, CacheValue{e.Value.Value, e.Value.Version})
		}
	}
}

func (n *Node) bufferize(req CacheEntity) {
	req.Value.Version = req.Value.Version.touch()
	n.buffer.Add(req)
	if n.bufferSize() > bufferSizeLimit {
		n.bufferSizeExceeded <- true
	}
	// fmt.Println("added to buffer", n.buffer.Buffer.Count())
}
