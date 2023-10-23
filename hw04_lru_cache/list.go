package hw04lrucache

type List interface {
	Len() int
	Front() *Node
	Back() *Node
	PushFront(v interface{}) *Node
	PushBack(v interface{}) *Node
	Remove(i *Node)
	MoveToFront(i *Node)
}

type Node struct {
	Value interface{}
	Next  *Node
	Prev  *Node
}

type list struct {
	head *Node
	tail *Node

	length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *Node {
	return l.head
}

func (l *list) Back() *Node {
	return l.tail
}

func (l *list) PushFront(v interface{}) *Node {
	newHead := &Node{
		Value: v,
		Next:  l.head,
	}

	if l.head != nil {
		l.head.Prev = newHead
	}

	if l.tail == nil {
		l.tail = newHead
	}

	l.head = newHead
	l.length++
	return newHead
}

func (l *list) PushBack(v interface{}) *Node {
	newTail := &Node{
		Value: v,
		Prev:  l.tail,
	}

	if l.tail != nil {
		l.tail.Next = newTail
	}

	if l.head == nil {
		l.head = newTail
	}

	l.tail = newTail
	l.length++
	return newTail
}

func (l *list) Remove(i *Node) {
	if i == l.head { // i.prev == nil
		l.head = i.Next
	}

	if i == l.tail { // i.next == nil
		l.tail = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *Node) {
	l.Remove(i)
	l.PushFront(i.Value)
}
