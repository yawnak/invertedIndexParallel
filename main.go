package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/asstronom/invertedIndexParallel/pkg/dict"
	"github.com/asstronom/invertedIndexParallel/pkg/index"
	"github.com/asstronom/invertedIndexParallel/pkg/server"
	"github.com/gin-gonic/gin"
)

var (
	numOfMappers  int
	numOfReducers int
)

func RunGin(idx *index.Index, filenamemap map[int]string) {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "Hello. To lookup words use endpoint /lookup/:word")
	})

	router.GET("/lookup/", func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "Hello. To lookup words use endpoint /lookup/:word")
	})

	router.GET("/lookup/:word", func(ctx *gin.Context) {
		word := ctx.Param("word")
		postingsList := idx.GetPostingsList(word)
		if postingsList == nil {
			ctx.String(http.StatusNotFound, "word not found")
			return
		}
		type Pl struct {
			Filename string
			Count    int
		}
		pls := make([]Pl, 0, len(postingsList.Postings))
		for _, v := range postingsList.Postings {
			pls = append(pls, Pl{Filename: filenamemap[int(v.Docid)], Count: int(v.Count)})
		}
		ctx.IndentedJSON(http.StatusOK, pls)
	})

	router.Run(":8080")
}

func main() {
	runtime.GOMAXPROCS(4)
	gin.SetMode(gin.ReleaseMode)
	flag.IntVar(&numOfMappers, "m", -1, "specify number of mappers")
	flag.IntVar(&numOfReducers, "r", -1, "specify number of reducers")
	flag.Parse()

	if numOfMappers == -1 {
		panic("number of mappers is not specified")
	}
	if numOfReducers == -1 {
		panic("number of reducers is not specified")
	}

	idx := index.NewIndex(numOfMappers, numOfReducers)
	dir, err := os.ReadDir("data")
	if err != nil {
		log.Fatalf("error opening dir: %s", err)
	}
	files := make([]io.Reader, 0, len(dir))
	filenamemap := make(map[int]string, len(dir))
	filenamedict := dict.NewDictionary(50)
	for i, fentry := range dir {
		//fmt.Println("data/" + fentry.Name())
		file, err := os.Open("data/" + fentry.Name())
		filenamemap[i] = fentry.Name()
		filenamedict.Insert(strconv.Itoa(i), fentry.Name())
		if err != nil {
			log.Fatalf("error opening file: %s", err)
		}
		files = append(files, file)
	}
	start := time.Now()
	idx.IndexDocs(files)
	elapsed := time.Since(start)
	fmt.Printf("Time to build the index: %s\n", elapsed.String())

	srv := server.NewServer(idx, filenamedict)
	go srv.Listen()
	//wg := sync.WaitGroup{}
	// wg.Add(1)
	// wg.Wait()
	RunGin(idx, filenamemap)
}
