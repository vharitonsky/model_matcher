package main

import (
    "log"
	"testing"
)

//type Product struct {
//    id, category_id, model_id, name string
//}

func TestMatchProducts(t *testing.T) {
	products := []Product{
		Product{id: "10", name: "5abc ahead", category_id: "180510"},
		Product{id: "11", name: "No category", category_id: "1234"},
		Product{id: "12", name: "No match", category_id: "12345"},
	}
    matched_products := MatchProducts(&products)
    log.Print(matched_products.Len())
	return
}
