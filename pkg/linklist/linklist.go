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
