package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	var node *ListItem

	// если список пустой
	if l.head == nil {
		node = &ListItem{Value: v}
		l.head = node
		l.tail = node
	} else {
		// вставка в начало
		node = &ListItem{Value: v, Next: l.head}
		l.head.Prev = node
		l.head = node
	}
	l.len++
	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	var node *ListItem

	// если список пустой
	if l.tail == nil {
		node := &ListItem{Value: v}
		l.head = node
		l.tail = node
	} else {
		// вставка в конец
		node := &ListItem{Value: v, Prev: l.tail}
		l.tail.Next = node
		l.tail = node
	}
	l.len++
	return node
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	i.Prev = nil
	i.Next = nil

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	next := i.Next
	prev := i.Prev

	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}

	// перемещаем элемент в начало
	i.Next = l.head
	l.head.Prev = i
	i.Prev = nil
	l.head = i
}
