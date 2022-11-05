package domain

import (
	"io"
)

type WordToken struct {
	Term  string
	Count int64
	Docid int64
}

type Posting struct {
	Docid int64
	Count int64
}

type PostingsList struct {
	Term    string
	Posting Posting
}

type FileToken struct {
	DocID int64
	File  io.Reader
}
