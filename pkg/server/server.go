package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/asstronom/invertedIndexParallel/pkg/dict"
	"github.com/asstronom/invertedIndexParallel/pkg/domain"
	"github.com/asstronom/invertedIndexParallel/pkg/index"
)

type Server struct {
	index     *index.Index
	filenames *dict.Dictionary
}

func NewServer(index *index.Index, filenames *dict.Dictionary) *Server {
	if index == nil {
		panic("index is nil")
	}

	return &Server{index: index, filenames: filenames}
}

func MakeCode(code int64) []byte {
	c := make([]byte, 8)
	binary.PutVarint(c, 500)
	return c
}

func ReadInt64(r io.Reader) (int64, error) {
	ba := make([]byte, 8)
	n, err := r.Read(ba)
	if err != nil {
		return 0, fmt.Errorf("error reading bytes: %w", err)
	}
	if n != 8 {
		return 0, fmt.Errorf("not enough bytes")
	}
	res, _ := binary.Varint(ba)
	return res, nil
}

func ReadBytes(r io.Reader, n int64) ([]byte, error) {
	ba := make([]byte, n)
	count, err := r.Read(ba)
	if err != nil {
		return nil, fmt.Errorf("error reading bytes: %w", err)
	}
	if int64(count) != n {
		return nil, fmt.Errorf("not enough bytes")
	}
	return ba, nil
}

func readRequest(r io.Reader) (*Request, error) {
	var req Request
	var err error
	req.Length, err = ReadInt64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading length of request: %w", err)
	}
	b, err := ReadBytes(r, req.Length)
	if err != nil {
		return nil, fmt.Errorf("error reading word of request: %w", err)
	}
	req.Word = string(b)
	return &req, nil
}
func (srv *Server) Handle(c net.Conn) {
	rd := bufio.NewReader(c)
	wr := bufio.NewWriter(c)
	for {
		word, ok, err := rd.ReadLine()
		if !ok {
			code := make([]byte, 8)
			binary.PutVarint(code, 400)
			wr.Write(code)
			continue
		}
		if err != nil {
			log.Printf("connection terminated with %s: %s", c.LocalAddr().String(), err)
			return
		}
		pl := srv.index.GetPostingsList(string(word))
		if pl == nil {
			code := make([]byte, 8)
			binary.PutVarint(code, 404)
			wr.Write(code)
			continue
		}

		res := make([]domain.PostingWName, len(pl.Postings))

		for i := range pl.Postings {
			fname, ok := srv.filenames.Get(strconv.FormatInt(pl.Postings[i].Docid, 10))
			if !ok {
				code := make([]byte, 8)
				binary.PutVarint(code, 404)
				wr.Write(code)
				continue
			}
			res[i].Filename = fname.(string)
			res[i].Count = pl.Postings[i].Count
		}
		code := make([]byte, 8)
		binary.PutVarint(code, 200)
		wr.Write(code)
		resjson, err := json.Marshal(res)
		if err != nil {
			code := make([]byte, 8)
			binary.PutVarint(code, 500)
			wr.Write(code)
		}
		wr.Write(resjson)
	}
}

func (srv *Server) Listen(ctx context.Context) error {
	l, err := net.Listen("tcp4", "8000")
	defer l.Close()
	if err != nil {
		return fmt.Errorf("error listening: %w", err)
	}
	for c, err := l.Accept(); ; {
		if err != nil {
			log.Printf("error accepting connection: %s", err)
		}
		go srv.Handle(c)
	}
}
