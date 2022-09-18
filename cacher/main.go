package cacher

import (
	"os"
	"fmt"
	"net"
	"net/rpc"
	"log"
	"sync"
	"dbcache/types"
	"dbcache/cluster"
)

type Node struct {
	n       *cluster.Node
	cache   map[string]string
	mu      sync.RWMutex
}

var thisNode *Node

type Req types.Req

type Resp types.Resp

func (n *Node) Initialize() {
	thisNode = n
	n.cache = make(map[string]string)
}

func (r *Req) Get(req Req, resp *Resp) error {
	thisNode.mu.RLock()
	val, ok := thisNode.cache[req.Key]
	thisNode.mu.RUnlock()
	*resp = Resp{ok, req.Key, val}
	return nil
}

func (r *Req) Put(req Req, resp *Resp) error {
	thisNode.mu.Lock()
	thisNode.cache[req.Key] = req.Value
	thisNode.mu.Unlock()
	*resp = Resp{true, req.Key, ""}
	return nil
}
