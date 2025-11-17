package utils

import (
	"errors"
)

const (
	NotExists = -1
)

type orderNode[T Comparable] struct {
	item *T
	next *orderNode[T]
}

type OrderedLinkedList[T Comparable] struct {
	start *orderNode[T]

	count int
}

func (l *OrderedLinkedList[T]) GetCount() int {
	return l.count
}

func (l *OrderedLinkedList[T]) Add(item *T) {
	l.count++

	if l.start == nil {
		l.start = &orderNode[T]{item: item, next: nil}
		return
	}

	var node = l.start
	if node.next == nil {
		var compare, err = (*node.item).CompareValues(*item)
		if err != nil {
			l.count--
			panic(err)
		}

		switch compare {
		case LESS:
			l.start = &orderNode[T]{item: item, next: node}
			return
		}
	}

	for node.next != nil {
		var compare, err = (*node.next.item).CompareValues(*item)
		if err != nil {
			panic(err)
		}

		if compare == LESS {
			node.next = &orderNode[T]{item: item, next: node.next}
			return
		}
		node = node.next
	}

	node.next = &orderNode[T]{item: item, next: nil}
	return
}

func (l *OrderedLinkedList[T]) Remove(item *T) error {
	if l.start == nil {
		return errors.New("list is empty")
	}

	if item == nil {
		return errors.New("item is nil")
	}

	var exact, err = (*item).ExactItem(*l.start.item)
	if err != nil {
		panic(err)
	}
	if exact {
		l.start = l.start.next
		l.count--
		return nil
	}

	var node = l.start

	for node.next != nil {
		exact, err = (*item).ExactItem(*node.next.item)
		if err != nil {
			panic(err)
		}

		if exact {
			node.next = node.next.next
			l.count--
			return nil
		}
		node = node.next
	}

	exact, err = (*item).ExactItem(*node.item)
	if err != nil {
		panic(err)
	}

	if exact {
		node.next = nil
		l.count--
		return nil
	}

	return errors.New("item not found")
}

// Contains return position of found element if it exists, otherwise -1
func (l *OrderedLinkedList[T]) Contains(item T) (int, error) {
	if l.start == nil {
		return 0, errors.New("list is empty")
	}

	var node = l.start
	var position = 0

	for node != nil {
		var exact, err = (*node.item).ExactItem(item)
		if err != nil {
			return 0, err
		}

		if exact {
			return position, nil
		}
		node = node.next
		position++
	}

	return NotExists, nil
}

func (l *OrderedLinkedList[T]) RemoveAt(index int) (*T, error) {
	if !(l.count > index) {
		return nil, errors.New("out of range")
	}

	if l.start == nil {
		return nil, errors.New("list is empty")
	}

	l.count--
	var node = l.start
	if index == 0 {
		l.start = l.start.next
		return node.item, nil
	}

	var i = 1
	for node.next != nil {
		if i == index {
			item := node.next.item
			node.next = node.next.next
			return item, nil
		}

		i++
		node = node.next
	}

	return nil, errors.New("item not found")
}
