package main

import (
	"testing"
)

func TestMatchProducts(t *testing.T) {
	products := []Product{
		Product{id: "10", name: "5abc ahead", category_id: "180510"},
		Product{id: "11", name: "No category", category_id: "1234"},
		Product{id: "12", name: "No match", category_id: "12345"},
	}
	matched_products := MatchProducts(&products)
	if !(matched_products.Len() == 1) {
		t.Error("At least one product should match")
	}
	matched_product := matched_products.Front().Value.(Product)
	if !(matched_product.id == "10") {
		t.Error("First product should have id 10, got", matched_product.id)
	}
	if !(matched_product.model_id == "703387") {
		t.Error("First product should have a model_id 703387 got", matched_product.model_id)
	}
	return
}
