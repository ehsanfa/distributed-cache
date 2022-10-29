package deque

import (
	"fmt"
	ll "github.com/ehsanfa/linked-list"
)

type Deque struct {
	list *ll.LinkedList
}

func (d *Deque) Count() int {
	return d.list.Count()
}

func (d *Deque) Enqueue(val interface{}) {
	d.list.Append(val)
}

func (d *Deque) Dequeue() (interface{}, error) {
	if d.Count() == 0 {
		return nil, fmt.Errorf("empty deque")
	}
	v := d.list.PopFirst()
	d.Enqueue(v.Value())
	return v.Value(), nil
}

func NewDeque() *Deque {
	d := new(Deque)
	d.list = &ll.LinkedList{}
	return d
}