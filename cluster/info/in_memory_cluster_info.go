package info

import (
	"dbcache/cluster/peer"
	"fmt"
	"log"
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

		// if !i.IsPeerAlive(peer) {
		// 	log.Println("UPDATE", peer, peerInfo)
		// 	i.Add(peer, peerInfo)
		// 	// peer.track(peerInfo)
		// 	// updatePartitionsInfo(peer, peerInfo)
		// 	continue
		// }

		pi, ok := i.getInfo(peer)
		if !ok || pi.Version().Number() < peerInfo.Version().Number() {
			log.Println("UPDATE", peer, peerInfo)
			i.Add(peer, peerInfo)
		}

		// updatePartitionsInfo(peer, peerInfo)

		// peer.track(peerInfo)

	}
}

func (i *InMemoryClusterInfo) GetClusterInfo() map[peer.Peer]peer.PeerInfo {
	return i.All()
}

func (i *InMemoryClusterInfo) Upsert(p peer.Peer, pi peer.PeerInfo) {
	i.Add(p, pi)
}

func (i *InMemoryClusterInfo) MarkAsDead(p peer.Peer) error {
	pi, ok := i.Get(p)
	if !ok {
		return fmt.Errorf("mark as dead failed. peer %s is unknown", p.Name())
	}
	pi = pi.MarkAsDead()
	log.Println("marked as dead", pi.IsAlive())
	i.Add(p, pi)
	log.Println("updated info", i.All())
	return nil
}

func (i *InMemoryClusterInfo) AllAlive() map[peer.Peer]peer.PeerInfo {
	l := make(map[peer.Peer]peer.PeerInfo)
	for p, pi := range i.All() {
		if !pi.IsAlive() {
			continue
		}
		l[p] = pi
	}
	return l
}
