package handlers

import (
	"log"
	"net/http"

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

	product := &data.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(product)
	w.WriteHeader(http.StatusCreated)
}

// func (p *Products) PutProduct(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Unable to convert id", http.StatusBadRequest)
// 		return
// 	}

// 	p.l.Println("Handle PUT Product", id)
// }
