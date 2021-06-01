package data

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product defines the structure for an API endpoint
// json
// validator
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format xxx-000-000
	re := regexp.MustCompile(`[a-zA-Z]+-[0-9]+-[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
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
