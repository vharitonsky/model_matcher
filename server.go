package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/vharitonsky/model_matcher/lib"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	modelsMap = make(map[string][]lib.Model)
	port      = flag.String("port", "8080", "port to run the server on")
)

type Product struct {
	Id          string `json:"id"`
	Category_id string `json:"category_id"`
	Model_id    string `json:"model_id"`
	Name        string `json:"name"`
}

type MatchData struct {
	Callback_url            string
	Callback_model_id_param string
	Products                []Product
}

func makeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func MatchProducts(products []Product) (matched_products []Product) {
	matched_products = make([]Product, 0)
	ch := make(chan interface{})
	for _, product := range products {
		product_name := lib.SplitName(lib.CleanName(product.Name))
		go func(product Product) {
			models, found := modelsMap[product.Category_id]
			if found {
				for _, model := range models {
					if lib.MatchNames(product_name, model.Name) {
						product.Model_id = model.Id
						ch <- product
						return
					}
				}
			}
			ch <- nil
		}(product)
	}
	for i := 0; i < len(products); i++ {
		matched_product := <-ch
		if matched_product != nil {
			matched_products = append(matched_products, matched_product.(Product))
		}
	}
	close(ch)
	return
}

func ProcessData(data []byte) (res []byte, callback_url string, err error) {
	var match_data MatchData
	err = json.Unmarshal(data, &match_data)
	if err != nil {
		return []byte{}, "", err
	}
	matched_products := MatchProducts(match_data.Products)
	if len(matched_products) > 0 {
		res, err = json.Marshal(matched_products)
		if err != nil {
			return []byte{}, "", err
		} else {
			return res, match_data.Callback_url, nil
		}
	} else {
		return []byte{}, "", errors.New("No products matched")
	}
}

func MatcherServer(w http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadAll(req.Body)
	log.Print("Received:" + string(data))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		res, callback_url, err := ProcessData(data)
		if err != nil {
			log.Fatal(err)
			http.Post(callback_url, "application/json", bytes.NewReader(res))
		}

	}()
	io.WriteString(w, "ok")
}

func init() {
	log.Print("Initializing models")
	models_count, categories_count := 0, 0
	start := time.Now()
	cat_file, err := os.Open("data/cats.txt")
	defer cat_file.Close()
	if err != nil {
		log.Fatal(err)
	}
	model_file_reader := bufio.NewReader(cat_file)
	for {
		cat_line, err := model_file_reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		cat_id := string(cat_line[:len(cat_line)-1])
		models := make([]lib.Model, 0)
		categories_count += 1

		file, err := os.Open("data/models/m_" + cat_id + ".txt")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(file)
		for {
			model_line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			parts := strings.Split(string(model_line[:len(model_line)-1]), "|")
			m := lib.Model{Id: parts[0], Name: lib.SplitName(parts[1])}
			models = append(models, m)
		}

		models_count += len(models)
		modelsMap[cat_id] = models
	}
	elapsed := time.Since(start)
	log.Print(fmt.Sprintf("Matcher initialized with %d models from %d categories in %s", models_count, categories_count, elapsed))
}

func main() {
	flag.Parse()
	log.Print("Running model matcher server on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, makeHandler(MatcherServer)))
}
