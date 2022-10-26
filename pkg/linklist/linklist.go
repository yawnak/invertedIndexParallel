package linklist

import "fmt"

type Node struct {
	next  *Node
	value interface{}
}

type LinkedList struct {
	Head     *Node
	capacity int
	length   int
}

func NewLinkedList(capacity int) *LinkedList {
	return &LinkedList{
		capacity: capacity,
	}
}
