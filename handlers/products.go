package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moewiz/go-microservice/data"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns the products from the data storage
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	// fetch the products from the data storage
	productList := data.GetProducts()

	// serialize the list to JSON
	err := productList.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
	w.WriteHeader(http.StatusCreated)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Cannot update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type KeyProduct struct{}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := data.Product{}

		if err := product.FromJSON(r.Body); err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		if err := product.Validate(); err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				w,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
