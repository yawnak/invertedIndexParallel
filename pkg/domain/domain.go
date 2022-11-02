package domain

import (
	"io"
)

type WordToken struct {
	Term  string
	Count int64
	Docid int64
}

type FileToken struct {
	DocID int64
	File  io.Reader
}
