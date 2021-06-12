package handlers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/moewiz/go-microservice/product-images/files"
)

// Files is a handler for reading and writing files
type Files struct {
	storage files.Storage
	log     *log.Logger
}

// Create a new Files handler
func NewFiles(s files.Storage, l *log.Logger) *Files {
	return &Files{storage: s, log: l}
}

// ServeHTTP implements the http.Handler interfacec
func (f *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Println("[INFO] Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, w, r)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, filename string, w http.ResponseWriter, r *http.Request) {
	f.log.Println("[INFO] Save file for product", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)
	err := f.storage.Save(fp, r.Body)
	if err != nil {
		f.log.Println("[ERROR] Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
