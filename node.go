package main

import (
	"os"
	"dbcache/cluster"
)


func main() {
	seeder_name, seeder_name_ok := os.LookupEnv("SEEDER_NAME")
	seeder_port, seeder_port_ok := os.LookupEnv("SEEDER_PORT")
	var n cluster.Node
	var s cluster.Seeder
	if seeder_name_ok && seeder_port_ok {
		s = cluster.CreateSeeder(seeder_name, seeder_port)
		n = cluster.CreateNode(false)
		n.SetSeeder(s)
	} else {
		n = cluster.CreateNode(true)
	}
	n.Initialize()
}