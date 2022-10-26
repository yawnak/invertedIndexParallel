package hash

func HashString(s string) int64 {
	var hash int64 = 7
	for char := range s {
		hash = hash*31 + int64(char)
	}
	return hash
}
