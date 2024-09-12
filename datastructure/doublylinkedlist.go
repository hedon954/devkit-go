package datastructure

// DoublyLinked represents the doubly linked list with dummy head and tail.
type DoublyLinked[T any] struct {
	count int
	head  *DoublyLinkedNode[T]
	tail  *DoublyLinkedNode[T]
}

// DoublyLinkedNode represents the node in the doubly linked list.
type DoublyLinkedNode[T any] struct {
	Value T
	prev  *DoublyLinkedNode[T]
	next  *DoublyLinkedNode[T]
}

func NewDoublyLinked[T any]() *DoublyLinked[T] {
	var zero T
	dummyHead := NewDoublyLinkedNode[T](zero)
	dummyTail := NewDoublyLinkedNode[T](zero)

	dummyHead.next = dummyTail
	dummyTail.prev = dummyHead

	return &DoublyLinked[T]{
		count: 0,
		head:  dummyHead,
		tail:  dummyTail,
	}
}

func NewDoublyLinkedNode[T any](value T) *DoublyLinkedNode[T] {
	return &DoublyLinkedNode[T]{
		Value: value,
		prev:  nil,
		next:  nil,
	}
}

func (d *DoublyLinked[T]) AddToHead(value T) *DoublyLinkedNode[T] {
	node := NewDoublyLinkedNode[T](value)
	node.prev = d.head
	node.next = d.head.next
	d.head.next.prev = node
	d.head.next = node

	d.count++
	return node
}

func (d *DoublyLinked[T]) AddToTail(value T) *DoublyLinkedNode[T] {
	node := NewDoublyLinkedNode[T](value)
	node.prev = d.tail.prev
	node.next = d.tail
	d.tail.prev.next = node
	d.tail.prev = node

	d.count++
	return node
}

func (d *DoublyLinked[T]) Remove(node *DoublyLinkedNode[T]) *DoublyLinkedNode[T] {
	if node == nil || node == d.head || node == d.tail || node.prev == nil || node.next == nil {
		return node
	}
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
	d.count--

	return node
}

func (d *DoublyLinked[T]) RemoveFromHead() *DoublyLinkedNode[T] {
	return d.Remove(d.head.next)
}

func (d *DoublyLinked[T]) RemoveFromTail() *DoublyLinkedNode[T] {
	return d.Remove(d.tail.prev)
}

func (d *DoublyLinked[T]) MoveToHead(node *DoublyLinkedNode[T]) *DoublyLinkedNode[T] {
	v := node.Value
	d.Remove(node)
	return d.AddToHead(v)
}

func (d *DoublyLinked[T]) MoveToTail(node *DoublyLinkedNode[T]) *DoublyLinkedNode[T] {
	v := node.Value
	d.Remove(node)
	return d.AddToTail(v)
}

func (d *DoublyLinked[T]) Count() int {
	return d.count
}

func (d *DoublyLinked[T]) Head() *DoublyLinkedNode[T] {
	if d.IsEmpty() {
		return nil
	}
	return d.head.next
}

func (d *DoublyLinked[T]) Tail() *DoublyLinkedNode[T] {
	if d.IsEmpty() {
		return nil
	}
	return d.tail.prev
}

func (d *DoublyLinked[T]) IsEmpty() bool {
	return d.count == 0
}

func (d *DoublyLinked[T]) Range(fn func(T) bool) {
	cur := d.head.next
	for cur != d.tail {
		if !fn(cur.Value) {
			break
		}
		cur = cur.next
	}
}
