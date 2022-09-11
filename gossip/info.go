package gossip

import (
	"sync"
)

var info map[Peer]PeerInfo
var infoMutex sync.RWMutex

func getInfo(peer Peer) (PeerInfo, bool){
	infoMutex.RLock()
	v, ok := info[peer]
	infoMutex.RUnlock()
	return v, ok
}

func setInfo(peer Peer, i PeerInfo) {
	infoMutex.Lock()
	info[peer] = i
	infoMutex.Unlock()
}