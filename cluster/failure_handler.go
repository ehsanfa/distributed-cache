package cluster

func (p *Peer) markAsDead() {
	i, ok := getInfo(*p)
	if ok {
		i.markAsDead()
		i.touch()
		setInfo(*p, i)
	}
}

func (n *Node) unbuddy(peer Peer) {
	peer.markAsDead()
	/**
	 * Find a better way to remove from the slice
	 */
	delete(n.buddies, peer)
}