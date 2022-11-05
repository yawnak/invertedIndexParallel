package index

import "github.com/asstronom/invertedIndexParallel/pkg/maps"

type Architect struct {
}

type Index struct {
	mappers []maps.Mapper
}

//n - number of mappers, m - numbers of reducers
func NewIndex(n int, m int) *Index {
	return &Index{mappers: make([]maps.Mapper, n)}
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

