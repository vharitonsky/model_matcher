package main

import (
	"bufio"
    "fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Model struct {
	id   string
	name string
}

func merge(cs []chan string) chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func readLines(file_path string) chan string {
	ch := make(chan string)
	go func() {
		file, err := os.Open(file_path)
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(file)
		defer file.Close()
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				break
			}
			ch <- strings.TrimSuffix(string(line), "\n")
		}
		close(ch)
	}()
	return ch
}

func splitLine(line string) (string, string) {
	parts := strings.Split(line, "|")
	return parts[0], parts[1]
}

func match_category(category_id string) chan string {
	ch := make(chan string)
	go func() {
		fmt.Println("Going to work on category_id:" + category_id)
		models := []Model{}

		for line := range readLines(fmt.Sprintf("models/m_%s.txt", category_id)) {
			id, name := splitLine(line)
			models = append(models, Model{id: id, name: name})
		}
		for line := range readLines(fmt.Sprintf("products/p_%s.txt", category_id)) {
			_, name := splitLine(line)
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
	for line := range readLines("cats.txt") {
		channels = append(channels, match_category(line))
	}
    l := len(channels)
	out := merge(channels)
	for _ = range out {
        count = count + 1
        fmt.Printf("\r %d/%d", count, l)
		//fmt.Println(msg)
	}
}
