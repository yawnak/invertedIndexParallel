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

func writeResponse(resp Response, wr io.Writer) error {
	bresp := make([]byte, 8+8, 8+8+resp.Length)

	binary.PutVarint(bresp[0:8], resp.Code)

	body, err := json.Marshal(resp.Body)
	if err != nil {
		return fmt.Errorf("error marshaling body: %w", err)
	}
	resp.Length = int64(len(body))

	binary.PutVarint(bresp[8:16], resp.Length)

	bresp = append(bresp, body...)
	_, err = wr.Write(bresp)
	if err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}
	return nil
}

func (srv *Server) Handle(c net.Conn) {
	rd := bufio.NewReader(c)
	defer c.Close()
	for {
		var resp Response
		req, err := readRequest(rd)
		if err != nil {
			err = writeResponse(Response{Code: 404}, c)
			if err != nil {
				fmt.Println(err)
			}
			break
		}
		fmt.Println(req)

		pl := srv.index.GetPostingsList(req.Word)
		if pl == nil {
			resp.Code = 404
		} else {
			pwn := make([]domain.PostingWName, len(pl.Postings))
			resp.Code = 200
			for i := range pwn {
				fname, ok := srv.filenames.Get(strconv.FormatInt(pl.Postings[i].Docid, 10))
				if !ok {
					resp.Code = 500
					break
				}
				pwn[i].Filename = fname.(string)
				pwn[i].Count = pl.Postings[i].Count
			}
			resp.Body = pwn
		}

		err = writeResponse(resp, c)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func (srv *Server) Listen() error {
	l, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Printf("error listening: %s\n", err)
		return fmt.Errorf("error listening: %w", err)
	}
	defer l.Close()
	fmt.Println("started waiting for connections!")
	for {
		c, err := l.Accept()
		log.Println("accepted connection!")
		if err != nil {
			log.Printf("error accepting connection: %s", err)
		}
		go srv.Handle(c)
	}
}
