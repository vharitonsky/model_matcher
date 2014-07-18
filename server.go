package main

import (
	"github.com/vharitonsky/goutil"
	//"github.com/vharitonsky/model_matcher/lib"
	"strings"
    "fmt"
)

var (
	modelsMap = make(map[string][]Model)
)

type Model struct{
    id, name string
}

//func match(name string, callback_url string) {

//}

func init() {
	for line := range(goutil.ReadLines("data/cats.txt")) {
		modelsMap[line] = make([]Model, 0)
		for model_line := range(goutil.ReadLines("data/models/m_" + line + ".txt")) {
			parts := strings.Split(model_line, "|")
            m := Model{id: parts[0], name: parts[1]}
			modelsMap[line] = append(modelsMap[line], m)
		}
	}
}

func main() {

}
