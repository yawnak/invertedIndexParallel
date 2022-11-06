package reduce

import (
	"sort"

	"github.com/asstronom/invertedIndexParallel/pkg/dict"
	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

type Reducer struct {
	d *dict.Dictionary
}

func NewReducer() *Reducer {
	return &Reducer{d: dict.NewDictionary(50)}
}

func (r *Reducer) Reduce(in <-chan domain.WordToken) {
	for tkn := range in {
		val, ok := r.d.Get(tkn.Term)
		var l domain.PostingsList
		if !ok {
			l = domain.PostingsList{
				Term: tkn.Term,
				Postings: []domain.Posting{{
					Docid: tkn.Docid,
					Count: tkn.Count,
				}},
			}
		} else {
			l = val.(domain.PostingsList)
			l.Postings = append(l.Postings, domain.Posting{Docid: tkn.Docid, Count: tkn.Count})
		}
		r.d.Insert(tkn.Term, l)
	}

	for kv := range r.d.Range() {
		pl := kv.Val.(domain.PostingsList)
		sort.Slice(pl.Postings, func(i, j int) bool {
			return pl.Postings[i].Docid < pl.Postings[j].Docid
		})
		r.d.Insert(kv.Key, pl)
	}
}

func (r *Reducer) GetPostingsList(term string) *domain.PostingsList {
	val, ok := r.d.Get(term)
	if !ok {
		return nil
	} else if val == nil {
		return nil
	}
	pl := val.(domain.PostingsList)
	return &pl
}
