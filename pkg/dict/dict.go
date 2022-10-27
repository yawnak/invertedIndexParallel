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

func (d *Dictionary) Get(key string) (interface{}, bool) {
	khash := hash.HashString(key)
	idx := khash % int64(len(d.buckets))
	bucklen := d.buckets[idx].GetLen()
	cur := d.buckets[idx].Head
	for i := 0; i < bucklen; i++ {
		if (cur.Value.(bucket).Hash == khash) && (cur.Value.(bucket).Key == key) {
			return cur.Value.(bucket).Val, true
		}
		cur = cur.Next
	}
	return nil, false
}

