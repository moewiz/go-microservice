package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/moewiz/go-microservice/product-api/data"
	"github.com/moewiz/go-microservice/product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "moewiz-service", log.LstdFlags)
	v := data.NewValidation()

	// Create the handlers
	productsHandler := handlers.NewProducts(l, v)

	// Create a new serve mux and register the handlers
	r := mux.NewRouter()

	getRouter := r.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", productsHandler.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.GetProduct)

	postRouter := r.Methods("POST").Subrouter()
	postRouter.HandleFunc("/products", productsHandler.CreateProduct)
	postRouter.Use(productsHandler.MiddlewareValidateProduct)

	putRouter := r.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.UpdateProduct)
	putRouter.Use(productsHandler.MiddlewareValidateProduct)

	deleteRouter := r.Methods("DELETE").Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productsHandler.DeleteProduct)

	// Serve API Docs page
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	getRouter.Handle("/README", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// Create a server
	server := &http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      ch(r),             // default handlers
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read the request from the client
		WriteTimeout: 10 * time.Second,  // max time to write the response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap signal and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until signal is received
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
