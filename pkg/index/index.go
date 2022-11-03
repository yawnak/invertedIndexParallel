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


