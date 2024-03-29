package gossip

// import (
// 	"dbcache/cluster/cluster"
// 	"dbcache/cluster/gossip/buddy"
// 	"dbcache/cluster/info"
// 	"time"
// )

// const gossipInterval = 5 * time.Second
// const gossipTimeout = 5 * time.Second

// type SeederBasedGossip struct {
// 	interval time.Duration
// 	timeout  time.Duration
// 	cluster  cluster.Cluster
// 	buddies  buddy.HasBuddies
// 	info     info.ClusterInfo
// }

// /*
//   - When a node gets initialized, it makes a request to its seeder to get its info
//   - Seeder sends its list of known nodes and the initialized node picks one to be buddy with
//   - If the seeder doesn't have a buddy yet itself, it becomes the node's buddy
//   - Otherwise it spreads the word that this nodes needs a buddy
//   - If the messages get to a node that could get the newjoiner as its buddy, it makes a
//     request to the newjoiner's BecomeBuddy function and becomes it buddy
//   - There should be max limit on how many buddies a node can have
// */
// func (g SeederBasedGossip) Initialize(c cluster.Cluster, i info.ClusterInfo, buddies buddy.HasBuddies) error {
// 	g.cluster = c
// 	g.buddies = buddies
// 	g.info = i
// 	timer := time.NewTicker(g.interval)
// 	go func() {
// 		for {
// 			select {
// 			case <-timer.C:
// 				g.spawn()
// 			}
// 		}
// 	}()
// 	return nil
// }

// func (g SeederBasedGossip) spawn() {
// 	for peer, _ := range g.buddies.AllBuddies() {
// 		go g.doGossip(peer)
// 	}
// }

// func (g SeederBasedGossip) doGossip(p peer.Peer) error {
// 	fmt.Println("goroutine counter", runtime.NumGoroutine())
// 	/**
// 	 * REFACTOR !!
// 	 * Also, please be less harsh. Give them more chance!
// 	 */
// 	// fmt.Println("doing gossip with ", n.getBuddies(), p)
// 	i, err := g.cluster.GetClusterInfo(p)
// 	if err != nil {
// 		// n.unbuddy(p)
// 		// fmt.Println("unbuddiying from ", p)
// 		return err
// 	}
// 	g.info.Update(i)

// 	// var resp Response
// 	// gossipRequest := GossipRequest{info}
// 	// timer := time.NewTimer(gossipTimeout)
// 	// call := c.Go("Node.Gossip", gossipRequest, &resp, nil)
// 	// select {
// 	// case <-call.Done:
// 	// 	timer.Stop()
// 	// case <-timer.C:
// 	// 	n.unbuddy(p)
// 	// 	return nil
// 	// }

// 	// updateInfo(resp)
// 	// n.checkForBuddies()

// 	return nil
// }

// type GossipRequest struct {
// 	Info map[Peer]PeerInfo
// }

// func (g GossipRequest) GetInfo() map[Peer]PeerInfo {
// 	return g.Info
// }

// func (n *Node) peersToGossip() map[Peer]bool {
// 	peers := make(map[Peer]bool)
// 	if !n.hasBuddy() {
// 		if n.isSeeder {
// 			return peers
// 		}
// 		peers[n.getSeeder().getPeer()] = true
// 	} else {
// 		peers = n.getBuddies()
// 	}
// 	return peers
// }

/*
 - When a node gets initialized, it makes a request to its seeder to get its info
 - Seeder sends its list of known nodes and the initialized node picks one to be buddy with
 - If the seeder doesn't have a buddy yet itself, it becomes the node's buddy
 - Otherwise it spreads the word that this nodes needs a buddy
 - If the messages get to a node that could get the newjoiner as its buddy, it makes a
   request to the newjoiner's BecomeBuddy function and becomes it buddy
 - There should be max limit on how many buddies a node can have
*/

// func (n *Node) startGossiping(endSignal <-chan bool) {
// 	timer := time.NewTicker(gossipInterval)
// 	go func() {
// 		for {
// 			select {
// 			case <-timer.C:
// 				n.spawnToGossip()
// 			case <-endSignal:
// 				break
// 			}
// 		}
// 	}()
// }

// var alreadyPushed bool

// // func (n *Node) doGossip(p Peer) error {
// // 	fmt.Println("goroutine counter", runtime.NumGoroutine())
// // 	/**
// // 	 * REFACTOR !!
// // 	 * Also, please be less harsh. Give them more chance!
// // 	 */
// // 	// fmt.Println("doing gossip with ", n.getBuddies(), p)
// // 	c, err := n.getConnection(p)
// // 	if err != nil {
// // 		n.unbuddy(p)
// // 		fmt.Println("unbuddiying from ", p)
// // 		return err
// // 	}

// // 	var resp Response
// // 	gossipRequest := GossipRequest{info}
// // 	timer := time.NewTimer(gossipTimeout)
// // 	call := c.Go("Node.Gossip", gossipRequest, &resp, nil)
// // 	select {
// // 	case <-call.Done:
// // 		timer.Stop()
// // 	case <-timer.C:
// // 		n.unbuddy(p)
// // 		return nil
// // 	}

// // 	updateInfo(resp)
// // 	n.checkForBuddies()

// // 	return nil
// // }

// func updateInfo(g GossipMaterial) {
// 	/**
// 	 * This is a mess. Come back later and fix it
// 	 */
// 	for peer, peerInfo := range g.GetInfo() {

// 		if _, ok := getInfo(peer); !ok {
// 			peer.track(peerInfo)
// 			updatePartitionsInfo(peer, peerInfo)
// 			continue
// 		}

// 		pi, _ := getInfo(peer)
// 		if pi.Version.compare(peerInfo.Version) >= 0 {
// 			// The message has nothing to give to us. Moving on
// 			continue
// 		}

// 		fmt.Println("SURPRISE. UPDATING", peer, pi, peerInfo)

// 		updatePartitionsInfo(peer, peerInfo)

// 		peer.track(peerInfo)

// 	}
// }

// func updatePartitionsInfo(peer Peer, peerInfo PeerInfo) {
// 	if peer.isAlive() {
// 		addToParitionsInfo(peer, peerInfo)
// 	}
// }

// func (n *Node) Gossip(req GossipRequest, resp *Response) error {
// 	updateInfo(req)
// 	*resp = Response{
// 		Info: info,
// 	}
// 	return nil
// }
