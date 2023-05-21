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

	client_handler "github.com/prestonchoate/go-commerce-catalog/Handlers/Clients"
	product_handler "github.com/prestonchoate/go-commerce-catalog/Handlers/Products"
)

const DEFAULT_PORT = "5000"

func main() {
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
		r.Post("/", product_handler.HandlePostProducts)
		r.Route("/{productID}", func(r chi.Router) {
			r.Use(product_handler.ProductsCtx)
			r.Get("/", product_handler.HandleGetProduct)
			r.Put("/", product_handler.HandlePutProduct)
			r.Delete("/", product_handler.HandleDeleteProduct)
		})
	})

	r.Route("/clients", func(r chi.Router) {
		r.Get("/", client_handler.HandleGetAllProducts)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Printf("Starting server on port: %v", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), r)

	if (errors.Is(err, http.ErrServerClosed)) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		log.Fatalf("error starting server %s\n", err)
	}
}
