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
	size    int
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

func (d *Dictionary) Insert(key string, val interface{}) {
	khash := hash.HashString(key)
	idx := khash % int64(len(d.buckets))
	bucklen := d.buckets[idx].GetLen()
	cur := d.buckets[idx].Head
	var isExists bool
	for i := 0; i < bucklen; i++ {
		if (cur.Value.(bucket).Hash == khash) && (cur.Value.(bucket).Key == key) {
			isExists = true
			break
		}
		cur = cur.Next
	}
	if isExists {
		cur.Value = bucket{
			Hash: khash,
			Key:  key,
			Val:  val,
		}
	} else {
		err := d.buckets[idx].Insert(bucket{
			Hash: khash,
			Key:  key,
			Val:  val,
		})
		if err != nil {
			log.Fatalln("need to do resize")
		} else {
			d.size++
		}
	}
}
