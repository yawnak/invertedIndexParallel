package domain

import (
	"io"
)

type WordToken struct {
	Term  string
	Docid int64
	Count int64
}

type Posting struct {
	Docid int64
	Count int64
}

type PostingsList struct {
	Term     string
	Postings []Posting
}

type FileToken struct {
	DocID int64
	File  io.Reader
}
