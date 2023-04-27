package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	product_handler "github.com/prestonchoate/go-commerce-catalog/Handlers/Products"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load from .env")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/products", product_handler.HandleProducts)
	port := os.Getenv("PORT")

	log.Printf("Starting server on port: %v", port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), mux)

	if (errors.Is(err, http.ErrServerClosed)) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server %s\n", err)
		os.Exit(1)
	}
}
