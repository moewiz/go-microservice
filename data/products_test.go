package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := Product{
		Name:  "moewiz",
		Price: 1.00,
		SKU:   "AOE-12345-123456",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
