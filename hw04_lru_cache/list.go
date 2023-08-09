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

	// Если список пустой
	if l.head == nil {
		node = &ListItem{Value: v}
		l.head = node
		l.tail = node
	} else {
		// Вставка в начало
		node = &ListItem{Value: v, Next: l.head}
		l.head.Prev = node
		l.head = node
	}
	l.len++
	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	var node *ListItem

	// Если список пустой
	if l.tail == nil {
		node := &ListItem{Value: v}
		l.head = node
		l.tail = node
	} else {
		// Вставка в конец
		node := &ListItem{Value: v, Prev: l.tail}
		l.tail.Next = node
		l.tail = node
	}
	l.len++
	return node
}

func (l *list) Remove(i *ListItem) {
	// Если элемент первый в списке
	if i.Prev == nil {
		if i.Next != nil {
			l.head = i.Next
		}
		l.head.Prev = nil
	} else {
		i.Prev.Next = i.Next
	}

	// Если элемент последний в списке
	if i.Next == nil {
		if i.Prev != nil {
			l.tail = i.Prev
		}
		l.tail.Next = nil
	} else {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = nil

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	// Если элемент 1 в списке - ничего не делаем
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
	//Определяем последний элемент списка
	l.tail = i.Prev

	// Перемещаем элемент в начало
	i.Next = l.head
	l.head.Prev = i
	i.Prev = nil
	l.head = i
}
