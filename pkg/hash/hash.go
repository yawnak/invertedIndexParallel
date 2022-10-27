package hash

import "hash/fnv"

func HashString(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
