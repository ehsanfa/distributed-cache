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
}

type Peer struct{
	Name      string
	Port      Port
}

type Node struct {
	isSeeder     bool
	Name         string
	Port         Port
	buddy        Buddy
	seeder       Seeder
	noBuddyPeers map[Peer]bool
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

func (n *Node) noBuddySlice() []Peer {
	var ps []Peer
	for peer, _ := range n.noBuddyPeers {
		ps = append(ps, peer)
	}
	return ps
}

type Response struct {
	Info      map[Peer]PeerInfo
	BuddyLook []Peer
	Version   Version
}

func (resp Response) GetInfo() map[Peer]PeerInfo {
	return resp.Info
}

func (resp Response) GetVersion() Version {
	return resp.Version
}

func (resp Response) GetBuddyLook() map[Peer]PeerInfo  {
	return resp.BuddyLook
}

type PotentialBuddies []Peer
type BuddyRequestResp struct {
	Res bool
}

type GossipMaterial interface {
	GetInfo()      map[Peer]PeerInfo
	GetVersion()   Version
	GetBuddyLook() map[Peer]PeerInfo 
}