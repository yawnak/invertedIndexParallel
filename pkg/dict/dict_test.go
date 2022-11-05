package dict

import (
	"bufio"
	"fmt"
	"os"
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

func TestDictionary(t *testing.T) {
	dir, err := os.ReadDir("test_data")
	if err != nil {
		t.Errorf("error reading dir: %s", err)
	}

	d := NewDictionary(50)
	m := make(map[string]interface{}, 50)

	for i := range dir {
		fmt.Println(dir[i].Name())
		file, err := os.Open("test_data/" + dir[i].Name())
		if err != nil {
			t.Errorf("error opening file: %s", err)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := string(scanner.Bytes())
			val, ok := d.Get(word)
			if !ok {
				d.Insert(word, 1)
			} else {
				n := val.(int)
				n++
				d.Insert(word, n)
			}
			val, ok = m[word]
			if !ok {
				m[word] = 1
			} else {
				if val == nil {
					fmt.Println(word, val, ok)
				}
				n := val.(int)
				n++
				m[word] = n
			}
		}
	}

	for kv := range d.Range() {
		val, ok := m[kv.Key]
		if !ok {
			t.Errorf("missing key in map: %#v", kv)
			continue
		}
		nd := kv.Val.(int)
		nm := val.(int)
		if nd != nm {
			t.Errorf("wrong val, %d != %d", nd, nm)
		}
	}
	if d.Len() != len(m) {
		t.Errorf("len is not correct: %d != %d", d.Len(), len(m))
	}
}

func TestRange(t *testing.T) {
	dict := NewDictionary(0)
	type kv struct {
		Key   string
		Value int
	}
	m := map[string]interface{}{}
	kvs := []kv{
		{"pudge", 13},
		{"judge", 15},
		{"sanandreas", 115},
		{"pudge", 121},
	}
	for _, v := range kvs {
		dict.Insert(v.Key, v.Value)
		m[v.Key] = v.Value
	}
	for kv := range dict.Range() {
		if kv.Val != m[kv.Key] {
			t.Errorf("wrong value of key: %s:%d != %s:%d", kv.Key, kv.Val, kv.Key, m[kv.Key])
		}
		fmt.Println(kv)
	}
}
