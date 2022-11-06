package maps

import (
	"bufio"
	"unicode"

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

func stripNonLetters(s string) string {
	word := []rune(s)
	startidx := 0
	endidx := 0
	for i := 0; i < len(word); i++ {
		if unicode.IsLetter(word[i]) {
			startidx = i
			break
		}
	}
	for i := len(word) - 1; i >= 0; i-- {
		if unicode.IsLetter(word[i]) {
			endidx = i + 1
			break
		}
	}
	return string(word[startidx:endidx])
}

func (m *Mapper) Map(filetokens []domain.FileToken, out chan<- domain.WordToken) {
	for _, filetoken := range filetokens {
		d := dict.NewDictionary(50)
		scanner := bufio.NewScanner(filetoken.File)
		if filetoken.File == nil {
			panic("filetoken is nil")
		}
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Bytes()
			word = []byte(stripNonLetters(string(word)))
			if len(word) == 0 {
				continue
			}
			tkn, ok := d.Get(string(word))
			if !ok {
				tkn = domain.WordToken{Docid: filetoken.DocID, Term: string(word), Count: 1}
				d.Insert(string(word), tkn)
			} else {
				posting := tkn.(domain.WordToken)
				posting.Count++
				d.Insert(string(word), posting)
			}
		}
		for tkn := range d.Range() {
			if tkn.Val == nil {
				continue
			}
			out <- tkn.Val.(domain.WordToken)
		}
	}
	close(out)
}
