package dict

import (
	"fmt"
	"testing"
)

func TestDictSimple(t *testing.T) {
	dict := NewDictionary(0)
	dict.Insert("pudge", 13)
	res, ok := dict.Get("pudge")
	if !ok {
		t.Errorf(`error, no key "pudge" in the dictionary\n`)
	}
	if res != 13 {
		t.Errorf(`error, res not valid, %d != %d`, res, 13)
	}
	res, ok = dict.Get("judge")
	if ok {
		t.Errorf(`error, got OK but "judge" is not in the dictionary`)
	}
	if res != nil {
		t.Errorf(`error, res != nil but "judge" is not in the dictionary`)
	}

}

func TestDict(t *testing.T) {
	dict := NewDictionary(0)
	type kv struct {
		Key   string
		Value int
	}
	kvs := []kv{
		{"pudge", 13},
		{"judge", 15},
		{"sanandreas", 115},
		{"pudge", 121},
	}
	for _, v := range kvs {
		dict.Insert(v.Key, v.Value)
	}
	for i := 1; i < len(kvs); i++ {
		v := kvs[i]
		res, ok := dict.Get(v.Key)
		fmt.Println(res, ok)
		if !ok {
			t.Errorf("!ok, when must be ok, key = %s", v.Key)
		}
		if res != v.Value {
			t.Errorf("wrong value, key %s, val %d != %d", v.Key, res, v.Value)
		}
	}
	fmt.Printf("size: %d\n", dict.size)
	if dict.size != 3 {
		t.Errorf("dict size is wrong %d != %d", dict.size, 3)
	}
}


func TestRange(t *testing.T) {
	dict := NewDictionary(0)
	type kv struct {
		Key   string
		Value int
	}
	kvs := []kv{
		{"pudge", 13},
		{"judge", 15},
		{"sanandreas", 115},
		{"pudge", 121},
	}
	for _, v := range kvs {
		dict.Insert(v.Key, v.Value)
	}
	for kv := range dict.Range() {
		fmt.Println(kv)
	}
}
