package server

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/asstronom/invertedIndexParallel/pkg/dict"
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

func (srv *Server) Handle(c net.Conn) {
	rd := bufio.NewReader(c)
	wr := bufio.NewWriter(c)
	for {
		word, ok, err := rd.ReadLine()
		if !ok {
			code := make([]byte, 8)
			binary.PutVarint(code, 400)
			wr.Write(code)
		}
		if err != nil {
			log.Printf("connection terminated with %s: %w", c.LocalAddr().String(), err)
			return
		}
		pl := srv.index.GetPostingsList(string(word))
		if pl == nil {
			code := make([]byte, 8)
			binary.PutVarint(code, 404)
			wr.Write(code)
		}

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
			log.Println("error accepting connection: %w", err)
		}
		go srv.Handle(c)
	}
}
