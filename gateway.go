package main

import (
	"dbcache/cluster/buffer"
	"dbcache/cluster/cacher"
	"dbcache/cluster/info"
	"dbcache/cluster/node"
	"dbcache/cluster/peer"
	portHelper "dbcache/cluster/port"
	"os"
)

func main() {
	seeder_name, seeder_name_ok := os.LookupEnv("SEEDER_NAME")
	seeder_port, seeder_port_ok := os.LookupEnv("SEEDER_PORT")

	if !seeder_name_ok || !seeder_port_ok {
		panic("Missing seeder info")
	}

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

	sprt, err := portHelper.ConvertPort(seeder_port)
	if err != nil {
		panic(err)
	}
	seeder := peer.CreateLocalPeer(seeder_name, sprt)

	cache := cacher.CreateInMemoryCache()
	i := info.CreateInMemoryClusterInfo()
	buff := buffer.CreateInMemoryBuffer()

	n, err := node.CreateGatewayNode(p, seeder, cache, i, buff)
	if err != nil {
		panic(err)
	}

	n.Run()
}
