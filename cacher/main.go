package main

import (
	"os"
	// "fmt"
	// "net"
	// "net/rpc"
	// "log"
	"sync"
	"dbcache/types"
	"dbcache/cluster"
)

var cache = map[string]string{}
var mu sync.RWMutex

type Req types.Req

type Resp types.Resp

func (r *Req) Get(req Req, resp *Resp) error {
	mu.RLock()
	val, ok := cache[req.Key]
	mu.RUnlock()
	*resp = Resp{ok, req.Key, val}
	return nil
}

func (r *Req) Put(req Req, resp *Resp) error {
	mu.Lock()
	cache[req.Key] = req.Value
	mu.Unlock()
	*resp = Resp{true, req.Key, ""}
	return nil
}

func main() {
	// req := new(Req)
 //    rpc.Register(req)

	// listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer listener.Close()

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

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Print(err)
	// 		continue
	// 	}
	// 	log.Print(conn.RemoteAddr())
	// 	rpc.ServeConn(conn)
	// }
}
