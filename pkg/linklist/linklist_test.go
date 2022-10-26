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
	cur := ll.Head
	for i := 0; i < size; i++ {
		fmt.Printf("%v ", cur.value)
		if cur.value.(int) != i {
			t.Errorf("error value of node, %v != %d", cur.value, i)
		}
		cur = cur.next
	}
	fmt.Println()
}
