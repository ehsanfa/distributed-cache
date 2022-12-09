package deque

import (
	"fmt"
	"sync"
)

type Deque struct {
	sync.Mutex
	list *LinkedList
}

func (d Deque) Count() int {
	return d.list.Count()
}

func (d Deque) IsEmpty() bool {
	return d.list.Head() == nil || d.list.Tail() == nil
}

func (d Deque) Enqueue(val interface{}) {
	d.list.Append(val)
}

func (d Deque) Display() {
	d.list.Display()
}

func (d Deque) Dequeue() (interface{}, error) {
	if d.Count() == 0 {
		return nil, fmt.Errorf("empty deque")
	}
	// d.Lock()
	v, _ := d.list.PopFirst()
	d.Enqueue(v.Value())
	// d.Unlock()
	// if err != nil {
	// 	return nil, fmt.Errorf("empty deque")
	// }
	return v.Value(), nil
}

func NewDeque() *Deque {
	d := new(Deque)
	d.list = &LinkedList{}
	return d
}
