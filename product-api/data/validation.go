package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ValidationError wraps the validators FieldError
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' fialed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errors := []string{}
	for _, error := range v {
		errors = append(errors, error.Error())
	}
	return errors
}

// Validation contains
type Validation struct {
	validate *validator.Validate
}

// NewValidation created a new Validation type
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return &Validation{validate}
}

// Validate the item
// for more detail the returned error can be cast into a
// validator.ValidationErrors collection
func (v *Validation) Validate(i interface{}) ValidationErrors {
	var returnErrors []ValidationError

	if errors, ok := v.validate.Struct(i).(validator.ValidationErrors); ok {
		if errors != nil {
			for _, error := range errors {
				if fe, ok := error.(validator.FieldError); ok {
					// cast the FieldError into our ValidationError and append to the slice
					ve := ValidationError{fe}
					returnErrors = append(returnErrors, ve)
				}
			}
		}
	}

	return returnErrors
}

func validateSKU(fl validator.FieldLevel) bool {
	// SKU must be in the format xxx-000-000
	re := regexp.MustCompile(`[a-zA-Z]+-[0-9]+-[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}
