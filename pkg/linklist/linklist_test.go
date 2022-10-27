package linklist

import (
	"fmt"
	"testing"
)

func TestLinkedInsert(t *testing.T) {
	size := 4
	ll := NewLinkedList(size)
	for i := 0; i < size; i++ {
		err := ll.Insert(i)
		if err != nil {
			t.Errorf("error inserting: %v\n", err)
		}
	}
	err := ll.Insert(5)
	if err == nil {
		t.Errorf("no error when list is full\n")
	}
	cur := ll.Head
	for i := 0; i < size; i++ {
		fmt.Printf("%v ", cur.Value)
		if cur.Value.(int) != i {
			t.Errorf("error value of node, %v != %d", cur.Value, i)
		}
		cur = cur.Next
	}
	fmt.Println()
}
