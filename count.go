package main

import (
	"fmt"
	"github.com/vharitonsky/goutil"
)

func counter(category_id string, ch chan int) {
	model_count, product_count := 0, 0
	for _ = range goutil.ReadLines("data/models/m_" + category_id + ".txt") {
		model_count = model_count + 1
	}
	for _ = range goutil.ReadLines("data/products/p_" + category_id + ".txt") {
		product_count = product_count + 1
	}
	ch <- (model_count * product_count)
}

func main() {
	ops := 0
	ch := make(chan int)
	rout_no := 0

	for line := range goutil.ReadLines("data/cats.txt") {
		go counter(line, ch)
		rout_no = rout_no + 1
	}
	for i := 0; i < rout_no; i++ {
		ops_in_file := <-ch
		ops = ops + ops_in_file
	}
	close(ch)
	fmt.Printf("Total ops number: %d\n", ops)
}
