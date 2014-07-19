package main

import (
	"github.com/vharitonsky/goutil"
	"log"
    "fmt"
	"strings"
	"sync"
)

var (
	modelsMap = make(map[string][]Model)
)

type Model struct {
	id, name string
}

func match(name string, callback_url string) string {
	return "ok"
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

}
