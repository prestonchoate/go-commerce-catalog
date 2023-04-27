package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joho/godotenv"
	product_handler "github.com/prestonchoate/go-commerce-catalog/Handlers/Products"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load from .env")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Route("/products", func(r chi.Router) {
		r.Get("/", product_handler.HandleGetAllProducts)
		//To do a GET on a Product ID
		r.Route("/{productID}", func(r chi.Router) {
			r.Use(product_handler.ProductsCtx)
			r.Get("/", product_handler.HandleGetProduct)
			//r.Put("/", handler for PUT /products/123)
			//r.Delete("/", handler for DELETE /products/123)
		})
	})

	port := os.Getenv("PORT")

	log.Printf("Starting server on port: %v", port)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), r)

	if (errors.Is(err, http.ErrServerClosed)) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server %s\n", err)
		os.Exit(1)
	}
}
