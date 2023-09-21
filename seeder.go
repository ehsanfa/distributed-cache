package main

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/node"
	"dbcache/cluster/peer"
	portHelper "dbcache/cluster/port"
	"dbcache/cluster/version"
	"os"
)

func main() {
	var peerName string
	var err error
	peerName, ok := os.LookupEnv("NAME")
	if !ok {
		peerName, err = os.Hostname()
		if err != nil {
			panic(err)
		}
	}

	port, port_ok := os.LookupEnv("PORT")
	var p peer.Peer
	if port_ok {
		prt, err := portHelper.ConvertPort(port)
		if err != nil {
			panic(err)
		}
		p = peer.CreateLocalPeer(peerName, prt)
	} else {
		p = peer.CreateLocalPeer(peerName, 0)
	}

	ver := version.CreateGenClockVersion(1)
	peerInfo := peer.CreateSimplePeerInfo(peer.Seeder, ver, true)

	cache := cacher.CreateInMemoryCache()
	i := info.CreateInMemoryClusterInfo()
	buff := buffer.CreateInMemoryBuffer()

	n, err := node.CreateSeederNode(p, cache, i, buff)
	if err != nil {
		panic(err)
	}

	i.Add(p, peerInfo)

	n.Run()
}
