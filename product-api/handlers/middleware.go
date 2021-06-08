package handlers

import (
	"context"
	"net/http"

	"github.com/moewiz/go-microservice/product-api/data"
)

// MiddlewareValidateProduct validates the product in the request before calls next http.HandlerFunc
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		if err := data.FromJSON(&product, r.Body); err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		// validate the product
		if errs := p.v.Validate(product); len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)

			// return the validation messages as an array
			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)

			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
