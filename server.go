package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	"github.com/vharitonsky/model_matcher/lib"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	redis_conn *redis.Client
	modelsMap  = make(map[string][]lib.Model)
	port       = flag.String("port", "8080", "port to run the server on")
	sigc       = make(chan os.Signal, 1)
	version    = ""
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

func CheckVersion() {
	current_version, _ := redis_conn.Cmd("get", "_model_matcher_version").Str()
	if current_version == "" {
		log.Print("Current server version is empty, exiting...")
		os.Exit(1)
	}
	log.Print("Current version is ", current_version)
	if current_version != version {
		InitModels()
		version = current_version
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

func InitModels() {
	log.Print("Initializing models")
	models_count, categories_count := 0, 0
	start := time.Now()
	cats, _ := redis_conn.Cmd("lrange", "_model_matcher_cats", 0, -1).List()
	for _, cat_id := range cats {
		model_lines, _ := redis_conn.Cmd("lrange", "_model_matcher_cat_"+cat_id, 0, -1).List()
		models := make([]lib.Model, 0)
		categories_count += 1
		for _, model_line := range model_lines {
			parts := strings.Split(string(model_line), "|")
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
	signal.Notify(sigc,
		syscall.SIGHUP,
	)
	go func() {
		for {
			_ = <-sigc
			InitModels()
		}
	}()
	var err error
	redis_conn, err = redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer redis_conn.Close()
	CheckVersion()
	ticker := time.NewTicker(time.Duration(10) * time.Minute)
	go func() {
		for {
			_ = <-ticker.C
			CheckVersion()
		}

	}()
	log.Print("Running model matcher server on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, makeHandler(MatcherServer)))
}
