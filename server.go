package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/vharitonsky/goutil"
	"io"
	"io/ioutil"
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

func matchProducts(products *[]Product) (matched_products []Product) {
	var wg sync.WaitGroup
	wg.Add(len(*products))
	for _, product := range *products {
		go func() {
			for _, model := range modelsMap[product.category_id] {
				if model.name == product.name {
					log.Print("Product", product.name, "vs", model.name)
					product.model_id = model.id
					matched_products = append(matched_products, product)
					break
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return matched_products
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
		matched_products := matchProducts(&match_data.products)
		if len(matched_products) > 0 {
            marshalled, err := json.Marshal(matched_products)
            if err != nil{
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
