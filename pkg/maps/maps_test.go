package maps

import (
	"testing"

	"github.com/asstronom/invertedIndexParallel/pkg/domain"
)

func TestMap(t *testing.T) {
	m := Mapper{}
	out := make(chan domain.Token)
	
	m.Map([[]struct{int64 *os.File}{0, nil}], out)
}
