package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/moewiz/go-microservice/handlers"
)

func main() {
	l := log.New(os.Stdout, "moewiz-service", log.LstdFlags)

	// Create the handlers
	productsHandler := handlers.NewProducts(l)

	// Create a new serve mux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", productsHandler)

	// Create a server
	server := &http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sm,                // default handlers
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read the request from the client
		WriteTimeout: 10 * time.Second,  // max time to write the response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
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
