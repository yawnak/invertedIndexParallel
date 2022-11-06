package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asstronom/invertedIndexParallel/pkg/index"
	"github.com/gin-gonic/gin"
)

func main() {
	idx := index.NewIndex(8, 4)
	dir, err := os.ReadDir("data")
	if err != nil {
		log.Fatalf("error opening dir: %s", err)
	}
	files := make([]io.Reader, 0, len(dir))
	filenamemap := make(map[int]string, len(dir))
	for i, fentry := range dir {
		//fmt.Println("data/" + fentry.Name())
		file, err := os.Open("data/" + fentry.Name())
		filenamemap[i] = fentry.Name()
		if err != nil {
			log.Fatalf("error opening file: %s", err)
		}
		files = append(files, file)
	}
	start := time.Now()
	idx.IndexDocs(files)
	elapsed := time.Since(start)
	fmt.Printf("Time to build the index: %s\n", elapsed.String())

	fmt.Println(idx.GetPostingsList("again"))

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello. To lookup words use endpoint /lookup/:word")
	})

	router.GET("/lookup/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello. To lookup words use endpoint /lookup/:word")
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
		ctx.IndentedJSON(http.StatusFound, pls)
	})

	router.Run(":8080")
}
