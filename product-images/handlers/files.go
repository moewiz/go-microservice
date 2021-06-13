package handlers

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

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

// UploadREST implements the http.Handler interfacec
func (f *Files) UploadREST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Println("[INFO] Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, w, r.Body)
}

func (f *Files) UploadMultipart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(128 * 1024); err != nil {
		f.log.Println("[ERROR] Bad request", err)
		http.Error(w, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if _, err := strconv.Atoi(id); err != nil {
		f.log.Println("[ERROR] Bad request", err)
		http.Error(w, "Expected integer id", http.StatusBadRequest)
		return
	}

	mf, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Println("[ERROR] Bad request", err)
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}

	f.log.Println("[INFO] Process form for id:", id)
	f.saveFile(id, mh.Filename, w, mf)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, filename string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Println("[INFO] Save file for product", "id", id, "filename", filename)

	fp := filepath.Join(id, filename)
	err := f.storage.Save(fp, r)
	if err != nil {
		f.log.Println("[ERROR] Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
