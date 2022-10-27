package dict

import (
	"log"

	"github.com/asstronom/invertedIndexParallel/pkg/hash"
	"github.com/asstronom/invertedIndexParallel/pkg/linklist"
)

const (
	linkListCap = 8
)

type bucket struct {
	Hash int64
	Key  string
	Val  interface{}
}

type Dictionary struct {
	buckets []linklist.LinkedList
}

func NewDictionary(bucketsNum int) *Dictionary {
	if bucketsNum == 0 {
		bucketsNum = 10
	}
	dict := Dictionary{
		buckets: make([]linklist.LinkedList, bucketsNum),
	}
	for i := range dict.buckets {
		dict.buckets[i] = *linklist.NewLinkedList(linkListCap)
	}
	return &dict
}
