package info

import (
	"dbcache/cluster/peer"
	"sync"
)

type InMemoryClusterInfo struct {
	info map[peer.Peer]peer.PeerInfo
	mu   sync.RWMutex
}

func CreateInMemoryClusterInfo() *InMemoryClusterInfo {
	return &InMemoryClusterInfo{info: make(map[peer.Peer]peer.PeerInfo)}
}

func (i *InMemoryClusterInfo) IsPeerKnown(p peer.Peer) bool {
	if _, ok := i.getInfo(p); !ok {
		return false
	}
	return true
}

func (i *InMemoryClusterInfo) Get(p peer.Peer) (peer.PeerInfo, bool) {
	return i.getInfo(p)
}

func (i *InMemoryClusterInfo) getInfo(p peer.Peer) (peer.PeerInfo, bool) {
	i.mu.RLock()
	v, ok := i.info[p]
	i.mu.RUnlock()
	return v, ok
}

func (i *InMemoryClusterInfo) IsPeerAlive(p peer.Peer) bool {
	info, ok := i.getInfo(p)
	if !ok || !info.IsAlive() {
		return false
	}
	return true
}

func (i *InMemoryClusterInfo) All() map[peer.Peer]peer.PeerInfo {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.info
}

func (i *InMemoryClusterInfo) List() []peer.Peer {
	info := i.All()
	var peers []peer.Peer
	for p := range info {
		peers = append(peers, p)
	}
	return peers
}

func (i *InMemoryClusterInfo) Add(p peer.Peer, info peer.PeerInfo) {
	i.mu.Lock()
	i.info[p] = info
	i.mu.Unlock()
}

func (i *InMemoryClusterInfo) Remove(p peer.Peer) {
	i.mu.Lock()
	delete(i.info, p)
	i.mu.Unlock()
}

func (i *InMemoryClusterInfo) Replace(info map[peer.Peer]peer.PeerInfo) {
	i.mu.Lock()
	i.info = info
	i.mu.Unlock()
}

func (i *InMemoryClusterInfo) Update(info map[peer.Peer]peer.PeerInfo) {
	for peer, peerInfo := range info {

		if _, ok := i.getInfo(peer); !ok {
			i.Add(peer, peerInfo)
			// peer.track(peerInfo)
			// updatePartitionsInfo(peer, peerInfo)
			continue
		}

		pi, _ := i.getInfo(peer)
		if pi.Version().Number() < peerInfo.Version().Number() {
			i.Add(peer, peerInfo)
		}

		// updatePartitionsInfo(peer, peerInfo)

		// peer.track(peerInfo)

	}
}

func (i *InMemoryClusterInfo) GetClusterInfo() map[peer.Peer]peer.PeerInfo {
	return i.All()
}
