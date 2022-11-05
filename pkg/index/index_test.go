package index

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	f1, err := os.Open("test_data/0_2.txt")
	f2, err := os.Open("test_data/1_3.txt")
	f3, err := os.Open("test_data/2_3.txt")
	f4, err := os.Open("test_data/3_4.txt")
	if err != nil {
		t.Errorf("error opening file: %s", err)
	}
	idx := NewIndex(2, 2)
	idx.IndexDocs([]io.Reader{f1, f2, f3, f4})
	fmt.Println(idx.GetPostingsList("movie"))
}
