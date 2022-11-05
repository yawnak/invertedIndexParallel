package reduce

import (
	"github.com/asstronom/invertedIndexParallel/pkg/dict"
	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

type Reducer struct {
	d *dict.Dictionary
}

func NewReducer() *Reducer {
	return &Reducer{d: dict.NewDictionary(50)}
}

func (r *Reducer) Reduce(chan []domain.WordToken) {
	
}
