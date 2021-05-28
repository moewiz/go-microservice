package main

import (
	"log"
	"net/http"
)

func main() {
	health_check := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Check good!")
	}

	http.HandleFunc("/health-check", health_check)

	http.ListenAndServe(":9090", nil)
}
