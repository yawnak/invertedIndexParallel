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
