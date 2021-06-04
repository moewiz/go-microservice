package main

import (
	"fmt"
	"testing"

	"github.com/moewiz/go-microservice/sdk/client"
	"github.com/moewiz/go-microservice/sdk/client/products"
)

func TestOutClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	products, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(products)
}
