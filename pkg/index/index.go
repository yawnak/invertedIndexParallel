package index

import (
	"fmt"
	"io"
	"sync"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
	"github.com/asstronom/invertedIndexParallel/pkg/hash"
	"github.com/asstronom/invertedIndexParallel/pkg/maps"
	"github.com/asstronom/invertedIndexParallel/pkg/reduce"
)

type Architect struct {
}

type Index struct {
	mappers  []maps.Mapper
	reducers []*reduce.Reducer
}

// n - number of mappers, m - numbers of reducers
func NewIndex(n int, m int) *Index {
	reducers := make([]*reduce.Reducer, m)
	for i := range reducers {
		reducers[i] = reduce.NewReducer()
	}
	return &Index{mappers: make([]maps.Mapper, n), reducers: reducers}
}

func buildFanIn(cs []chan []domain.WordToken) <-chan []domain.WordToken {
	var wg sync.WaitGroup
	out := make(chan []domain.WordToken)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan []domain.WordToken) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func formFiletokens(files []io.Reader, startidx int) []domain.FileToken {
	res := make([]domain.FileToken, len(files))
	for i := range files {
		res[i] = domain.FileToken{
			DocID: int64(startidx) + int64(i),
			File:  files[i],
		}
	}
	return res
}

func (idx *Index) IndexDocs(files []io.Reader) {
	mapsout := make([]chan []domain.WordToken, len(idx.mappers))
	for i := range mapsout {
		mapsout[i] = make(chan []domain.WordToken)
	}
	fanin := buildFanIn(mapsout)
	pagesize := len(files) / len(idx.mappers)
	for i := 0; i < len(idx.mappers)-1; i++ {
		fts := formFiletokens(files[i*pagesize:(i+1)*pagesize], i*pagesize)
		fmt.Println(fts)
		go idx.mappers[i].Map(fts, mapsout[i])
	}
	fts := formFiletokens(files[(len(idx.mappers)-1)*pagesize:], (len(idx.mappers)-1)*pagesize)
	fmt.Println(fts)
	go idx.mappers[len(idx.mappers)-1].Map(fts, mapsout[len(idx.mappers)-1])

	reduceins := make([]chan []domain.WordToken, len(idx.reducers))
	for i := range reduceins {
		reduceins[i] = make(chan []domain.WordToken, 1)
	}
	for i := range idx.reducers {
		go idx.reducers[i].Reduce(reduceins[i])
	}
	for in := range fanin {
		h := hash.HashString(in[0].Term)
		reduceins[h%uint64(len(reduceins))] <- in
	}
	for i := range reduceins {
		close(reduceins[i])
	}
}
