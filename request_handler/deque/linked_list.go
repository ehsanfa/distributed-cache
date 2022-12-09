package deque

import (
	"fmt"
)

type Node struct {
	value interface{}
	next  *Node
	prev  *Node
}

func (n *Node) Value() interface{} {
	return n.value
}

type LinkedList struct {
	head  *Node
	tail  *Node
	count int
}

// adds a node to the head of the list
func (ll *LinkedList) Prepend(value interface{}) {
	// create a new node
	n := ll.newNode(value)

	//set the prev node of the new node as the old head
	oldhead := ll.head
	newhead := n
	if newhead != oldhead {
		newhead.prev = oldhead
	}

	//set the next node of the old head as the new head
	oldhead.next = newhead

	// set the head of the list as the new node
	ll.head = newhead
	if ll.isEmpty() {
		ll.tail = n
	}
	ll.count++
}

// adds a node to the tail of the list
func (ll *LinkedList) Append(value interface{}) {
	// create a new node
	n := ll.newNode(value)

	// set the next node of the new node as the old tail
	oldtail := ll.tail
	newtail := n
	if newtail != oldtail {
		newtail.next = oldtail
	}

	// set the prev node of the old tail as the new tail
	if oldtail != nil {
		oldtail.prev = newtail
	}

	// set the tail of the list as the new node
	ll.tail = newtail
	if ll.isEmpty() {
		ll.head = n
	}
	ll.count++
}

func (ll *LinkedList) newNode(value interface{}) *Node {
	n := &Node{}
	n.value = value
	if ll.isEmpty() {
		ll.head = n
		ll.tail = n
	}
	return n
}

func (ll *LinkedList) Pop() Node {
	n := Node{}
	n.value = ll.tail.value
	ll.Remove(ll.Tail())
	return n
}

func (ll *LinkedList) PopFirst() (Node, error) {
	n := Node{}
	if ll.head == nil {
		return n, fmt.Errorf("empty list")
	}
	n.value = ll.head.value
	ll.Remove(ll.Head())
	return n, nil

}

func (ll *LinkedList) isEmpty() bool {
	return ll.count == 0
}

func NewLinkedList() LinkedList {
	ll := LinkedList{}
	return ll
}

func (ll *LinkedList) Display() {
	list := ll.tail
	for list != nil {
		fmt.Printf("%+v ->", list.value)
		list = list.next
	}
	fmt.Println()
}

func (ll *LinkedList) Head() *Node {
	return ll.head
}

func (ll *LinkedList) Tail() *Node {
	return ll.tail
}

func (ll *LinkedList) Remove(n *Node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		ll.tail = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		ll.head = n.prev
	}
	n.prev = nil
	n.next = nil
	ll.count--
}

func (ll *LinkedList) Count() int {
	return ll.count
}

func (ll *LinkedList) RemoveValue(value interface{}) bool {
	n := ll.tail
	for n != nil {
		if value == n.value {
			ll.Remove(n)
			return true
		}
		n = n.next
	}
	return false
}

func (ll *LinkedList) Contains(value interface{}) bool {
	n := ll.tail
	for n != nil {
		if value == n.value {
			return true
		}
		n = n.next
	}
	return false
}
