package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Files is a handler for reading and writing files
type Files struct {
	storePath string
	log       *log.Logger
}

// Create a new Files handler
func NewFiles(s string, l *log.Logger) *Files {
	return &Files{storePath: s, log: l}
}

// ServeHTTP implements the http.Handler interfacec
func (f *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Println("[INFO] Handle POST", "id=", id, "filename=", filename)

	// check that the filepath is a valid name and file
	if id == "" || filename == "" {
		return
	}

	f.saveFile(id, filename, w, r)
}

func (f *Files) saveFile(id, filename string, w http.ResponseWriter, r *http.Request) {
	f.log.Println("[INFO] Saving file to physical disk")
}
