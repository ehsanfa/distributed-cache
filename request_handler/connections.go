package request_handler

import (
	"fmt"
	"net/rpc"
)

func (p *Peer) seederListen(connectionOpened chan<- bool) {
	if p.putChan == nil {
		p.putChan = make(chan putReq, 20)
	}
	var err error
	p.conn, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
	if err != nil {
		// fmt.Println(err, node)
		fmt.Println(err)
		connectionOpened <- false
		return
	}
	connectionOpened <- true
	defer p.conn.Close()
	for {
		select {
		case v := <-p.putChan:
			req := CacheRequest{Action: 2, Key: v.key, Value: v.val}
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			// fmt.Println(resp)
		case req := <-p.reqChan:
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			// fmt.Println("resp for put", resp)
		}
	}
}

func (p *Peer) listen() {
	defer p.conn.Close()
	for {
		select {
		case v := <-p.putChan:
			req := CacheRequest{Action: 2, Key: v.key, Value: v.val}
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			// fmt.Println(resp)
		case req := <-p.reqChan:
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			// fmt.Println("resp for put", resp)
		}
	}
}

func (p *Peer) req() {
	for {
		select {
		case req := <-p.reqChan:
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
			// fmt.Println(resp)
		}
	}
}

func dial(p *Peer) (*rpc.Client, error) {
	// fmt.Println(p.Name, p.Port)
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.Name, p.Port))
}