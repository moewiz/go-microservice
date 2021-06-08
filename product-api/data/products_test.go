package data

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameAndSKU(t *testing.T) {
	p := Product{
		Price: 1.99,
	}

	v := NewValidation()
	err := v.Validate(p)

	assert.Len(t, err, 2)
	fmt.Println(err)
}

func TestProductWrongPrice(t *testing.T) {
	p := Product{
		Name:  "moewiz",
		Price: -1.99,
		SKU:   "AOE-12345-123456",
	}

	v := NewValidation()
	err := v.Validate(p)

	assert.Len(t, err, 1)
	fmt.Println(err)
}

func TestProductInvalidSKU(t *testing.T) {
	p := Product{
		Name:  "moewiz",
		Price: 1.99,
		SKU:   "AOE-ABC-XYZ",
	}

	v := NewValidation()
	err := v.Validate(p)

	assert.Len(t, err, 1)
	fmt.Println(err)
}

func TestValidProduct(t *testing.T) {
	p := Product{
		Name:  "moewiz",
		Price: 1.99,
		SKU:   "AOE-12345-123456",
	}

	v := NewValidation()
	err := v.Validate(p)

	assert.Len(t, err, 0)
	fmt.Println(err)
}

func TestProductsToJSON(t *testing.T) {
	lp := Products{
		{Name: "Coffee"}, {Name: "Tea"},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(lp, b)

	assert.NoError(t, err)
}
