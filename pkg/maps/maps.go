package maps

import (
	"bufio"
	"fmt"

	"github.com/asstronom/invertedIndexParallel/pkg/dict"
	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

type Mapper struct {
}

func linSearch(docid int64, tokens []domain.WordToken) int {
	for i := range tokens {
		if docid == tokens[i].Docid {
			return i
		}
	}
	return -1
}

func (m *Mapper) Map(filetokens []domain.FileToken, out chan<- []domain.WordToken) {
	d := dict.NewDictionary(50)
	for _, filetoken := range filetokens {
		scanner := bufio.NewScanner(filetoken.File)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Bytes()
			tkn, ok := d.Get(string(word))
			if !ok {
				tkn = []domain.WordToken{{Docid: filetoken.DocID, Term: string(word), Count: 1}}
				d.Insert(string(word), tkn)
			} else {
				postings_list := tkn.([]domain.WordToken)
				idx := linSearch(filetoken.DocID, postings_list)
				if idx == -1 {
					postings_list = append(postings_list, domain.WordToken{Docid: filetoken.DocID, Term: string(word), Count: 1})
				} else {
					postings_list[idx].Count++
				}
				d.Insert(string(word), postings_list)
			}
		}
	}
	fmt.Println("map")
	for tkn := range d.Range() {
		if tkn.Val == nil {
			continue
		}
		out <- tkn.Val.([]domain.WordToken)
	}
	close(out)
}
