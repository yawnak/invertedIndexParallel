package domain

import "os"

type WordToken struct {
	term  string
	count int64
	docid int64
}

type FileToken struct {
	DocID int64
	File  *os.File
}
