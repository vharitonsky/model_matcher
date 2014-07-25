package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vharitonsky/goutil"
	"github.com/vharitonsky/model_matcher/lib"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	modelsMap = make(map[string]*list.List)
	port      = flag.String("port", "8080", "port to run the server on")
)

type Product struct {
	id, category_id, model_id, name string
}

type MatchData struct {
	callback_url            string
	callback_model_id_param string
	products                []Product
}

func makeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func MatchProducts(products *[]Product) (matched_products *list.List) {
	var wg sync.WaitGroup
	matched_products = list.New()
	wg.Add(len(*products))
	var model lib.Model
	for _, product := range *products {
		product_name := lib.SplitName(lib.CleanName(product.name))
		go func(product Product) {
			l, found := modelsMap[product.category_id]
			if found {
				for e := l.Front(); e != nil; e = e.Next() {
					model = e.Value.(lib.Model)
					if lib.MatchNames(product_name, model.Name) {
						product.model_id = model.Id
						matched_products.PushBack(product)
						break
					}
				}
			}
			wg.Done()
		}(product)
	}
	wg.Wait()
	return
}

func MatcherServer(w http.ResponseWriter, req *http.Request) {
	var match_data MatchData
	data, err := ioutil.ReadAll(req.Body)
	log.Print("Received:" + string(data))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		json.Unmarshal(data, &match_data)
		matched_products := MatchProducts(&match_data.products)
		if matched_products.Len() > 0 {
			marshalled, err := json.Marshal(matched_products)
			if err != nil {
				log.Fatal(err)
			}
			http.Post(match_data.callback_url, "application/json", bytes.NewReader(marshalled))
		}
	}()

	io.WriteString(w, "ok")
}

func init() {
	log.Print("Initializing models")
	var wg sync.WaitGroup
	models_count, categories_count := 0, 0
	for line := range goutil.ReadLines("data/cats.txt") {
		modelsMap[line] = list.New()
		categories_count += 1
		wg.Add(1)
		go func(cat_id string) {
			for model_line := range goutil.ReadLines("data/models/m_" + cat_id + ".txt") {
				parts := strings.Split(model_line, "|")
				m := lib.Model{Id: parts[0], Name: lib.SplitName(parts[1])}
				modelsMap[cat_id].PushBack(m)
				models_count += 1
			}
			wg.Done()
		}(line)
	}
	wg.Wait()
	log.Print(fmt.Sprintf("Matcher initialized with %d models from %d categories", models_count, categories_count))
}

func main() {
	flag.Parse()
	log.Print("Running model matcher server on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, makeHandler(MatcherServer)))
}
