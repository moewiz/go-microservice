package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/moewiz/go-microservice/data"
)

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ErrorInvalidProductPath is an error message when the product path is not valid
var ErrorInvalidProductPath = fmt.Errorf("invalid path, path should be /products/{id}")

// Generic Error is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// swagger:route GET /products products listProducts
// Returns a list of products
//
// Responses:
//	200: productsResponse
//	500: errorResponse

// GetProducts returns the products from the data storage
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] Handle GET products")
	// fetch the products from the data storage
	productList := data.GetProducts()

	// serialize the list to JSON
	err := data.ToJSON(productList, w)
	if err != nil {
		http.Error(w, "[ERROR] Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products listSingle
// Return a single product from the database
//
// Responses:
//	200: productResponse
//	404: errorResponse
//	500: errorResponse

// GetProduct returns the product by product ID
func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := getProductID(r)
	p.l.Printf("[DEBUG] Handle GET product/%d", productID)

	product, err := data.GetProduct(productID)
	if err == data.ErrorProductNotFound {
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Internal server error"}, w)
		return
	}
	err = data.ToJSON(product, w)
	if err != nil {
		// we should never be here but log the error just in case
		http.Error(w, "[ERROR] Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /products products createProduct
// Create a new product
//
// Responses:
//	201: productResponse
//	400: errorValidation
//	422: errorValidation
//	500: errorResponse

func (p *Products) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", product)

	data.AddProduct(&product)
	w.WriteHeader(http.StatusCreated)
}

// swagger:route PUT /products/{id} products updateProduct
// Update the product
//
//	Responses:
//		204: noContentResponse
//		404: errorResponse
//		500: errorResponse
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Printf("Handle PUT product/%d", id)

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err := data.UpdateProduct(id, &product)
	if err == data.ErrorProductNotFound {
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Internal server error"}, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Delete the product from the database
//
//	Responses:
//		204: noContentResponse
//		404: errorResponse
//		500: errorResponse
func (p *Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := getProductID(r)
	p.l.Printf("[DEBUG] Handle DELETE product/%d", productID)

	err := data.DeleteProduct(productID)
	if err == data.ErrorProductNotFound {
		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Internal server error"}, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getProductID returns the product ID from the URL
// Panics if CANNOT convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse url
	vars := mux.Vars(r)

	// convert id to integer and return
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen but panic error just in case
		panic(err)
	}

	return productID
}
