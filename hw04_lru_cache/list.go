package hw04lrucache

type List interface {
	Len() uint64 // длина списка не может быть со знаком и битность лучше определить
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
	len        uint64
	head, tail *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() uint64 {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.head,
	}
	if l.head == nil {
		l.tail = newItem
	} else {
		l.head.Prev = newItem
	}
	l.head = newItem
	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.tail,
	}
	if l.tail == nil {
		l.head = newItem
	} else {
		l.tail.Next = newItem
	}
	l.tail = newItem
	l.len++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if l.len == 0 {
		return
	}
	if i == l.head {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i == l.tail {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	i.Prev.Next = i.Next
	i.Prev = nil
	l.head.Prev = i
	i.Next = l.head
	l.head = i
}
