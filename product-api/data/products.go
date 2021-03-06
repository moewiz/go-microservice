package data

import (
	"fmt"
	"time"
)

// Product defines the structure for an API endpoint
// swagger:model
type Product struct {
	// the id for this product
	//
	// required: true
	// unique: true
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-zA-Z]+-[0-9]+-[0-9]+
	SKU string `json:"sku" validate:"required,sku"`

	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func GetProduct(productID int) (*Product, error) {
	position := findProduct(productID)

	if position == -1 {
		return nil, ErrorProductNotFound
	}

	return productList[position], nil
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	position := findProduct(id)

	if position == -1 {
		return ErrorProductNotFound
	}

	p.ID = id
	productList[position] = p
	return nil
}

func DeleteProduct(id int) error {
	position := findProduct(id)

	if position == -1 {
		return ErrorProductNotFound
	}

	productList = append(productList[:position], productList[position+1:]...)

	return nil
}

var ErrorProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffe",
		Price:       29000,
		SKU:         "MW1",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       25000,
		SKU:         "MW1",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}
