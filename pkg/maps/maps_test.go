package maps

import (
	"fmt"
	"os"
	"testing"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

func TestMap(t *testing.T) {
	res := map[string]domain.WordToken{
		"Once":  {Term: "Once", Docid: 0, Count: 2},
		"again": {Term: "again", Docid: 0, Count: 1},
		"yes":   {Term: "yes", Docid: 0, Count: 1},
		"no":    {Term: "no", Docid: 0, Count: 1},
	}

	m := Mapper{}
	f1, err := os.Open("test_data/0_2.txt")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}
	sl := make([]domain.FileToken, 1)
	sl[0] = domain.FileToken{DocID: 0, File: f1}
	out := make(chan domain.WordToken)
	go m.Map(sl, out)
	for tkn := range out {
		if val, ok := res[tkn.Term]; !ok {
			t.Errorf("odd value %#v\n", val)
		}
		if res[tkn.Term] != tkn {
			t.Errorf("wrong value: %#v != %#v", tkn, res[tkn.Term])
		}
		fmt.Println(tkn)
	}
}

func TestStrip(t *testing.T) {
	ss := []string{
		"movie", "movie...", "movies...", "i'm", "i'm.", "!i'm",
	}
	assert := []string{
		"movie", "movie", "movies", "i'm", "i'm", "i'm",
	}
	for i := range ss {
		if stripNonLetters(ss[i]) != assert[i] {
			t.Errorf("wrong result: %s != %s", stripNonLetters(ss[i]), assert[i])
		}
		fmt.Println(stripNonLetters(ss[i]))
	}
}
