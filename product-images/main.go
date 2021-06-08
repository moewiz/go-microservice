package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/moewiz/go-microservice/product-images/config"
)

func main() {
	l := log.New(os.Stdout, "product-images", log.LstdFlags)

	if err := godotenv.Load(); err != nil {
		l.Fatal("[ERROR] Error loading .env file")
	}
	conf := config.NewConfig()
	bindAddress := fmt.Sprintf("%s:%d", conf.Server.BindAddress, conf.Server.PORT)

	// Create a new serve mux and register handlers
	sm := mux.NewRouter()

	// Create a server
	server := &http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      sm,                // default handlers
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read the request from the client
		WriteTimeout: 10 * time.Second,  // max time to write the response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Start the server
	go func() {
		l.Println("[INFO] Starting the server:", bindAddress)
		if err := server.ListenAndServe(); err != nil {
			l.Println("[ERROR] Unable to start server", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("[INFO] Received terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
