package maps

import (
	"fmt"
	"os"
	"testing"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

func TestMap(t *testing.T) {
	m := Mapper{}
	f1, err := os.Open("test_data/0_2.txt")
	f2, err := os.Open("test_data/1_3.txt")
	f3, err := os.Open("test_data/2_3.txt")
	f4, err := os.Open("test_data/3_4.txt")
	if err != nil {
		t.Errorf("error opening file: %s", err)
	}
	sl := make([]domain.FileToken, 4)
	sl[0] = domain.FileToken{DocID: 0, File: f1}
	sl[1] = domain.FileToken{DocID: 1, File: f2}
	sl[2] = domain.FileToken{DocID: 2, File: f3}
	sl[3] = domain.FileToken{DocID: 3, File: f4}
	out := make(chan []domain.WordToken)
	go m.Map(sl, out)
	for tkn := range out {
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
