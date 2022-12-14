package cluster

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	partition "dbcache/partitioning"
)

type Port uint16

func (p Port) String() string {
	return fmt.Sprintf("%d", uint32(p))
}

type Seeder Peer

func (s Seeder) getPeer() Peer {
	return Peer(s)
}

type PeerInfo struct {
	Version   Version
	IsAlive   bool
	Partition partition.Partition
}

func NewPeerInfo() PeerInfo {
	v := Version{1, 1}
	return PeerInfo{Version: v, IsAlive: true}
}

func (p *PeerInfo) touch() {
	// mu.Lock()
	p.Version.touch()
	// mu.Unlock()
}

func (p *PeerInfo) markAsDead() {
	p.IsAlive = false
}

type Peer struct{
	Name      string
	Port      Port
}

func (p *Peer) isKnown() bool {
	if _, ok := getInfo(*p); !ok {
		return false
	}
	return true
}

func (p *Peer) isAlive() bool {
	if i, ok := getInfo(*p); !ok || !i.IsAlive {
		return false
	}
	return true
}

func (p *Peer) hasPartition(pi PeerInfo) bool {
	if pi.Partition == (partition.Partition{}) {
		return false
	}
	return true
}

func (p *Peer) track(i PeerInfo) {
	setInfo(*p, i)
}

func (p *Peer) setPort(port Port) {
	p.Port = port
}

func (p *Peer) setName(n string) {
	p.Name = n
}

func (p *Peer) isSame(other Peer) bool {
	return p.Name == other.Name && p.Port == other.Port
}

type CacheVersion int

type CacheValue struct {
	Value   string
	Version CacheVersion
}

type Node struct {
	bufferSizeExceeded chan bool
	cacheVersionsMu    sync.RWMutex
	cacheVersions      map[string]CacheVersion
	connections        map[Peer]*rpc.Client
	partitions         []partition.Partition    
	partition          partition.Partition
	isSeeder           bool
	bufferMu           sync.RWMutex
	buddies            map[Peer]bool
	cacheMu            sync.RWMutex
	seeder             Seeder
	buffer             Buffer
	cache              map[string]CacheValue
	Peer               *Peer
	mu                 sync.RWMutex
}

func (n *Node) SetSeeder(s Seeder) {
	n.seeder = s
}

func (n *Node) setPort(listener net.Listener) {
	port := Port(listener.Addr().(*net.TCPAddr).Port)
	n.getPeer().setPort(port)
}

func (n *Node) getPeer() *Peer {
	return n.Peer
}

func (n *Node) getSeeder() Seeder {
	return n.seeder
}

func (n *Node) getPeerInfo() PeerInfo {
	p := n.getPeer()
	var i PeerInfo
	if !p.isKnown() {
		i = NewPeerInfo()
	} else {
		i, _ = getInfo(*p)
	}
	return i
}

func (n *Node) setName(name string) {
	n.getPeer().setName(name)
}

type Response struct {
	Info      map[Peer]PeerInfo
	Cache     map[string]string
	Partition partition.Partition
}

func (resp Response) GetInfo() map[Peer]PeerInfo {
	return resp.Info
}

type BuddyRequestResp struct {
	Res bool
}

type GossipMaterial interface {
	GetInfo()      map[Peer]PeerInfo
}