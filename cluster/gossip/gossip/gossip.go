package gossip

import (
	innercircle "dbcache/cluster/gossip/inner_circle"
	"dbcache/cluster/info"
	"dbcache/cluster/network"
	"dbcache/cluster/peer"
	"fmt"
	"log"
	"time"
)

type Seeder peer.Peer

const maxBuddyNum = 3
const gossipInterval = 1 * time.Second
const minglingInterval = 1 * time.Second
const timeout = 5 * time.Second

type Gossip struct {
	network      network.Network
	info         info.ClusterInfo
	inner_cirlce innercircle.InnerCircle
}

func CreateGossipNetwork(
	network network.Network,
	info info.ClusterInfo,
	seeder peer.Peer,
	isSeeder bool,
	self peer.Peer,
) (*Gossip, error) {
	if seeder == nil && !isSeeder {
		return nil, fmt.Errorf("Seeder cannot be nil when the node itself is not a seeder")
	}
	inner_cirle := innercircle.CreateInMemoryBuddies(maxBuddyNum, self)
	if !isSeeder {
		if ok := inner_cirle.Add(seeder); !ok {
			return nil, fmt.Errorf("failed to add the seeder to the inner circle")
		}
	}
	return &Gossip{network, info, inner_cirle}, nil
}

func (g *Gossip) Start() {
	gossipTimer := time.NewTicker(gossipInterval)
	minglingInterval := time.NewTicker(minglingInterval)
	go func() {
		for {
			select {
			case <-gossipTimer.C:
				g.spawn()
			case <-minglingInterval.C:
				g.mingle()
			}
		}
	}()
}

func (g *Gossip) spawn() {
	for peer := range g.inner_cirlce.All() {
		log.Println("gossiping with ", peer)
		go g.gossip(peer)
	}
}

func (g *Gossip) gossip(p peer.Peer) {
	node, err := g.network.Connect(p, timeout)
	if err != nil {
		g.gossipFailed(p, err)
		return
	}
	info, err := node.GetClusterInfo()
	if err != nil {
		g.gossipFailed(p, err)
	} else {
		g.info.Update(info)
	}
}

func (g *Gossip) gossipFailed(p peer.Peer, err error) {
	g.info.MarkAsDead(p)
	g.inner_cirlce.Remove(p)
	log.Println(err)
}

func (g *Gossip) mingle() {
	info := g.info.AllAlive()
	candidates := make([]peer.Peer, 0)
	for p := range info {
		candidates = append(candidates, p)
	}
	log.Println("candidates to mingle", candidates)
	if len(info) == 0 {
		return
	}

	g.inner_cirlce.Shuffle(candidates)
}
