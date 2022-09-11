package gossip

import (
	"fmt"
	"net"
)

type Port uint16

func (p Port) String() string {
	return fmt.Sprintf("%d", uint32(p))
}

type Seeder Peer

type PeerInfo struct {
	Version   Version
	IsAlive   bool
	HasBuddy  bool
	IsSomeonesBuddy bool
}

func NewPeerInfo() PeerInfo {
	v := Version{1, 1}
	return PeerInfo{Version: v, IsAlive: true, HasBuddy: false}
}

func (p *PeerInfo) touch() {
	p.Version.touch()
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

func (p *Peer) track(i PeerInfo) {
	setInfo(*p, i)
}

type Node struct {
	isSeeder     bool
	Name         string
	Port         Port
	buddy        Buddy
	seeder       Seeder
	noBuddyPeers map[Peer]PeerInfo
	buddyWith    []Peer
	version      Version
}

func (n *Node) SetSeeder(s Seeder) {
	n.seeder = s
}

func (n *Node) setPort(listener net.Listener) {
	n.Port = Port(listener.Addr().(*net.TCPAddr).Port)
}

func (n *Node) getPeer() Peer {
	return Peer{n.Name, n.Port}
}

func (n *Node) getPeerInfo() PeerInfo {
	p := n.getPeer()
	var i PeerInfo
	if !p.isKnown() {
		i = NewPeerInfo()
		p.track(i)
	} else {
		i, _ = getInfo(p)
	}
	return i
}

func (n *Node) noBuddySlice() []Peer {
	var ps []Peer
	for peer, _ := range n.noBuddyPeers {
		ps = append(ps, peer)
	}
	return ps
}

type Response struct {
	Info      map[Peer]PeerInfo
}

func (resp Response) GetInfo() map[Peer]PeerInfo {
	return resp.Info
}

type PotentialBuddies []Peer
type BuddyRequestResp struct {
	Res bool
}

type GossipMaterial interface {
	GetInfo()      map[Peer]PeerInfo
}