package request_handler

import (
	"fmt"
	"net/rpc"
)

const putReqBuffer = 20

func (p *Peer) seederListen(connectionOpened chan<- bool) {
	if p.putChan == nil {
		p.putChan = make(chan putReq, 20)
	}
	var err error
	p.conn, err = rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.info.Name, p.info.Port))
	if err != nil {
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
		case req := <-p.reqChan:
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
		}
	}
}

func (p *Peer) prepare() error{
	conn, err := dial(p)
	if err != nil {
		return err
	}
	p.conn = conn
	p.reqChan = make(chan CacheRequest, putReqBuffer)
	return nil
}

func (p *Peer) listen() {
	for {
		select {
		case req := <-p.reqChan:
			var resp CacheRequestResponse
			p.conn.Call("Node.Put", req, &resp)
		}
	}
	p.conn.Close()
}

func dial(p *Peer) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", p.info.Name, p.info.Port))
}