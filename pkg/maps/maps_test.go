package maps

import (
	"fmt"
	"os"
	"testing"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

func TestMap(t *testing.T) {
	res := map[string][]domain.WordToken{
		"Once":  {{Term: "Once", Docid: 0, Count: 2}, {Term: "Once", Docid: 1, Count: 1}},
		"again": {{Term: "again", Docid: 0, Count: 1}, {Term: "again", Docid: 1, Count: 1}, {Term: "again", Docid: 2, Count: 1}, {Term: "again", Docid: 3, Count: 1}},
		"yes":   {{Term: "yes", Docid: 0, Count: 1}},
		"no":    {{Term: "no", Docid: 0, Count: 1}, {Term: "no", Docid: 3, Count: 1}},
		"I'm":   {{Term: "I'm", Docid: 1, Count: 1}},
		"you":   {{Term: "you", Docid: 1, Count: 1}},
		"It":    {{Term: "It", Docid: 2, Count: 1}},
		"is":    {{Term: "is", Docid: 2, Count: 1}},
		"going": {{Term: "going", Docid: 2, Count: 1}},
		"well":  {{Term: "well", Docid: 2, Count: 1}},
		"way":   {{Term: "way", Docid: 3, Count: 1}},
	}

	m := Mapper{}
	f1, err := os.Open("test_data/0_2.txt")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}
	f2, err := os.Open("test_data/1_3.txt")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}
	f3, err := os.Open("test_data/2_3.txt")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}
	f4, err := os.Open("test_data/3_4.txt")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}
	sl := make([]domain.FileToken, 4)
	sl[0] = domain.FileToken{DocID: 0, File: f1}
	sl[1] = domain.FileToken{DocID: 1, File: f2}
	sl[2] = domain.FileToken{DocID: 2, File: f3}
	sl[3] = domain.FileToken{DocID: 3, File: f4}
	out := make(chan []domain.WordToken)
	go m.Map(sl, out)
	for tkn := range out {
		assertvals := res[tkn[0].Term]
		if len(assertvals) != len(tkn) {
			t.Errorf("wrong assert len: %d != %d, key: %s", len(tkn), len(assertvals), tkn[0].Term)
		}
		for i := range tkn {
			if assertvals[i].Docid != tkn[i].Docid {
				t.Errorf("wrong WordToken: %#v != %#v", tkn[i], assertvals[i])
			}
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
