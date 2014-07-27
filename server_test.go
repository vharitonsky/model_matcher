package main

import (
	"testing"
)

func TestMatchProducts(t *testing.T) {
	products := []Product{
		Product{Id: "10", Name: "5abc ahead", Category_id: "180510"},
		Product{Id: "11", Name: "No category", Category_id: "1234"},
		Product{Id: "12", Name: "No match", Category_id: "12345"},
	}
	matched_products := MatchProducts(&products)
	if !(matched_products.Len() == 1) {
		t.Error("At least one product should match")
	}
	matched_product := matched_products.Front().Value.(Product)
	if !(matched_product.Id == "10") {
		t.Error("First product should have id 10, got", matched_product.Id)
	}
	if !(matched_product.Model_id == "703387") {
		t.Error("First product should have a model_id 703387 got", matched_product.Model_id)
	}
	return
}

func TestProcessData(t *testing.T) {
	data := []byte(`
        {"callback_url": "1234", "products": [{"Id": "10", "Name": "5abc ahead", "Category_id": "180510"}]}
    `)
	ProcessData(data)
	return
}
