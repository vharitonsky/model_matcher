package main

import (
    "fmt"
	"github.com/vharitonsky/goutil"
)

type Model struct {
	id   string
	name string
}

func match_category(category_id string) chan string {
	ch := make(chan string)
	go func() {
		fmt.Println("Going to work on category_id:" + category_id)
		models := []Model{}

		for line := range goutil.ReadLines(fmt.Sprintf("models/m_%s.txt", category_id)) {
			id, name := goutil.SplitLine(line)
			models = append(models, Model{id: id, name: name})
		}
		for line := range goutil.ReadLines(fmt.Sprintf("products/p_%s.txt", category_id)) {
			_, name := goutil.SplitLine(line)
			for _, model := range models {
				fmt.Printf("%s vs %s\n", name, model.name)
			}
		}
        ch <- "1"
        fmt.Println("Finished processing category_id:" + category_id)
		close(ch)
	}()
	return ch

}

func main() {
    count := 0
	channels := []chan string{}
	for line := range goutil.ReadLines("cats.txt") {
		channels = append(channels, match_category(line))
	}
    l := len(channels)
	out := goutil.MergeChannels(channels)
	for _ = range out {
        count = count + 1
        fmt.Printf("\r %d/%d", count, l)
		//fmt.Println(msg)
	}
}
