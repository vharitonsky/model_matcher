package main

import (
	"flag"
	"fmt"
	"github.com/vharitonsky/goutil"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	modelsMap = make(map[string][]Model)
	port      = flag.String("port", "8080", "port to run the server on")
)

type Model struct {
	id, name string
}

type Product struct {
	id, category_id, name string
}

type RequestData struct {
	callback_url            string
	callback_model_id_param string
	products                []Product
}

func makeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func MatcherServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello")
}

func init() {
	log.Print("Initializing models")
	var wg sync.WaitGroup
	models_count, categories_count := 0, 0
	for line := range goutil.ReadLines("data/cats.txt") {
		modelsMap[line] = make([]Model, 0)
		categories_count += 1
		wg.Add(1)
		go func() {
			for model_line := range goutil.ReadLines("data/models/m_" + line + ".txt") {
				parts := strings.Split(model_line, "|")
				m := Model{id: parts[0], name: parts[1]}
				modelsMap[line] = append(modelsMap[line], m)
				models_count += 1
			}
			wg.Done()
		}()
	}
	wg.Wait()
	log.Print(fmt.Sprintf("Matcher initialized with %d models from %d categories", models_count, categories_count))
}

func main() {
	flag.Parse()
	log.Print("Running model matcher server on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, makeHandler(MatcherServer)))
}
