package products

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/prestonchoate/go-commerce-catalog/Models"
	product_repository "github.com/prestonchoate/go-commerce-catalog/Models/Products"
)

func HandleProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling %v to %v", r.Method, r.RequestURI)
	if r.Method == http.MethodGet {
		HandleGetAllProducts(w, r)
	}
}

func HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	products, err := product_repository.GetAll()
	if err != nil {
		log.Print("Failed to retrieve products")
		resp["error"] = err.Error()
		json_resp, _ := generateResponse(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write(json_resp)
		return
	}
	resp["products"] = products
	json_resp, _ := generateResponse(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func handlePostProducts(w http.ResponseWriter, r *http.Request) {
	// validate body contains valid product struct
	// try to do db insert
	// return new product to client
}

func HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product, ok := ctx.Value("product").(*models.Product)
	if !ok {
		log.Print("Product not found in request context")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	resp := make(map[string]any)
	resp["product"] = product
	json_resp, err := generateResponse(resp)
	if err != nil {
		log.Printf("Could not convert product id %v into JSON", product.ID)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	original_product, ok := ctx.Value("product").(*models.Product)
	if !ok {
		log.Print("Product not found in request context")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
	}
	var request_product models.Product
	err := json.NewDecoder(r.Body).Decode(&request_product)
	if err != nil {
		log.Print("Could not parse request body as product")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	} 
	updated_product, err := product_repository.UpdateProductById(*original_product, request_product)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	resp := make(map[string]any)
	resp["product"] = updated_product
	json_resp, err := generateResponse(resp)
	if err != nil {
		log.Printf("Could not convert product id %v into JSON", updated_product.ID)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func generateResponse(resp map[string]any) ([]byte, error){
	json_resp, err := json.Marshal(resp)
	if (err != nil) {
		log.Printf("Error occured in JSON marshal. Err: %s\n", err)
		return nil, errors.New("could not marshal JSON response")
	}
	return json_resp, nil
}

func ProductsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product_id_param := chi.URLParam(r, "productID")
		product_id, err := strconv.Atoi(product_id_param)
		if err != nil {
			log.Printf("Could not convert %v to int", product_id_param)
			return
		}
		product, err := product_repository.GetProduct(product_id)
		if err != nil {
			log.Printf("Could not retrieve product ID: %v", product_id)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "product", &product)
		//printContextInternals(ctx, true)
		next.ServeHTTP(w, r.WithContext(ctx))
  })
}