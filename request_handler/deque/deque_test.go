package deque

import (
	"fmt"
	"testing"
)

type PeerInfo struct {
	Name   string
	Port   int
}

type Peer struct {
	info    PeerInfo
}

func TestDeque(t *testing.T) {
	var v interface{}
	var err error
	d := NewDeque()
	_, err = d.Dequeue()
	if err == nil {
		t.Error("failed dequing empty deque")
	}
	d.Enqueue("test1")
	d.Enqueue("test2")
	d.Enqueue("test3")
	if d.Count() != 3 {
		t.Error("miscount")
	}
	v, _ = d.Dequeue()
	if d.Count() != 3 {
		t.Error("enqueue failed")
	}
	if v != "test1" {
		t.Error("dequeue failed", v, "test1")
	}
	v, _ = d.Dequeue()
	if v != "test2" {
		t.Error("dequeue failed", v, "test2")
	}
	v, _ = d.Dequeue()
	if v != "test3" {
		t.Error("dequeue failed", v, "test3")
	}
	v, _ = d.Dequeue()
	if v != "test1" {
		t.Error("dequeue failed", v, "test1")
	}
}

type PeerInt interface {
	getPeer() *Peer
}

func (p *Peer) getPeer() *Peer {
	return p
}

func ShowPeer(p *Peer) {
	fmt.Println(p)
}

func TestPeer(t *testing.T) {
	p := &Peer{PeerInfo{"test", 1234}}
	d := NewDeque()
	d.Enqueue(p)
	v, _ := d.Dequeue()
	ShowPeer(v.(PeerInt).getPeer())
}