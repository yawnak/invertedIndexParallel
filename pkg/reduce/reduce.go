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

func (r *Reducer) Reduce(in <-chan []domain.WordToken) {
	for tkns := range in {
		for i := range tkns {
			val, ok := r.d.Get(tkns[i].Term)
			var l domain.PostingsList
			if !ok {
				l = domain.PostingsList{
					Term: tkns[i].Term,
					Postings: []domain.Posting{{
						Docid: tkns[i].Docid,
						Count: tkns[i].Count,
					}},
				}
			} else {
				l = val.(domain.PostingsList)
				l.Postings = append(l.Postings, domain.Posting{Docid: tkns[i].Docid, Count: tkns[i].Count})
			}
			r.d.Insert(tkns[i].Term, l)
		}
	}
}
