package reduce

import (
	"fmt"
	"sync"
	"testing"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

func TestReduce(t *testing.T) {
	input := []domain.WordToken{
		{Term: "Once", Docid: 0, Count: 2},
		{Term: "Once", Docid: 3, Count: 5},
		{Term: "again", Docid: 2, Count: 1}, {Term: "again", Docid: 3, Count: 1},
		{Term: "again", Docid: 0, Count: 1}, {Term: "again", Docid: 1, Count: 1},
	}

	assert := []domain.PostingsList{
		{Term: "Once", Postings: []domain.Posting{{Docid: 0, Count: 2}, {Docid: 3, Count: 5}}},
		{Term: "again", Postings: []domain.Posting{{Docid: 0, Count: 1}, {Docid: 1, Count: 1}, {Docid: 2, Count: 1}, {Docid: 3, Count: 1}}},
	}

	reducer := NewReducer()
	inchan := make(chan domain.WordToken)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		reducer.Reduce(inchan)
		wg.Done()
	}()
	for i := range input {
		inchan <- input[i]
	}
	close(inchan)
	wg.Wait()
	for i := range assert {
		postlist := reducer.GetPostingsList(assert[i].Term)
		fmt.Println(postlist)
		for j := range postlist.Postings {
			if postlist.Postings[j] != assert[i].Postings[j] {
				t.Errorf("wrong posting: %#v != %#v", postlist, assert[i])
			}
		}
	}

	postlist := reducer.GetPostingsList("no")
	if postlist != nil {
		t.Errorf("wrong posting, must be nil but %#v", postlist)
	}

}
