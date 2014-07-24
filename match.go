package main

import (
	"fmt"
	"github.com/vharitonsky/goutil"
)

type Model struct {
	id   string
	name string
}

func match_category(category_id string) chan interface{} {
	ch := make(chan interface{})
	go func() {
		count := 0
		fmt.Println("Going to work on category_id:" + category_id)
		models := []Model{}

		for line := range goutil.ReadLines(fmt.Sprintf("data/models/m_%s.txt", category_id)) {
			id, name := goutil.SplitLine(line)
			models = append(models, Model{id: id, name: name})
		}
		for line := range goutil.ReadLines(fmt.Sprintf("data/products/p_%s.txt", category_id)) {
			_, _ = goutil.SplitLine(line)
			for _, _ = range models {
				count = count + 1
				//ch <- fmt.Sprintf("%s vs %s\n", name, model.name)
			}
		}
		ch <- fmt.Sprintf(" %s: %d", category_id, count)
		close(ch)
	}()
	return ch

}

func main() {
	count := 0
	channels := make([]chan interface{}, 0)
	for line := range goutil.ReadLines("data/cats.txt") {
		channels = append(channels, match_category(line))
	}
	l := len(channels)
	out := goutil.MergeChannels(channels)
	for cnt := range out {
		count = count + 1
		fmt.Println(cnt)
		fmt.Printf("\r %d/%d", count, l)
	}
}
