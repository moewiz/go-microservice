package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/moewiz/go-microservice/product-images/config"
	"github.com/moewiz/go-microservice/product-images/files"
	"github.com/moewiz/go-microservice/product-images/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-images", log.LstdFlags)

	if err := godotenv.Load(); err != nil {
		l.Fatal("[ERROR] Error loading .env file")
		os.Exit(1)
	}
	conf := config.NewConfig()
	serverAddress := fmt.Sprintf("%s:%d", conf.Server.BindAddress, conf.Server.PORT)

	// create the storage class, use local storage
	// max filesize 5MB
	store, err1 := files.NewLocalStorage(conf.Storage.BasePath, 5*1000*1024)
	if err1 != nil {
		l.Println("[ERROR] Unable to create storage", "error", err1)
		os.Exit(1)
	}

	// Create the Files handler
	filesHandler := handlers.NewFiles(store, l)
	// Create a new serve mux and register handlers
	sm := mux.NewRouter()
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/images/{id:[0-9]+}/{filename:[a-z\\-A-Z]+\\.[a-z]{3}}", filesHandler.UploadREST)
	postRouter.HandleFunc("/", filesHandler.UploadMultipart)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.Handle(
		"/images/{id:[0-9]+}/{filename:[a-z\\-A-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(conf.Storage.BasePath))),
	)

	// Create a server
	server := &http.Server{
		Addr:         serverAddress,     // configure the bind address
		Handler:      ch(sm),            // default handlers
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read the request from the client
		WriteTimeout: 10 * time.Second,  // max time to write the response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Start the server
	go func() {
		l.Println("[INFO] Starting the server", serverAddress)
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
