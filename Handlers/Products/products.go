package products

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/prestonchoate/go-commerce-catalog/Models"
	product_repository "github.com/prestonchoate/go-commerce-catalog/Models/Products"
)

func HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	products, err := product_repository.GetAll()
	if err != nil {
		log.Print("Failed to retrieve products")
		resp["error"] = err.Error()
		writeResponse(resp, w, http.StatusBadGateway)
		return
	}
	resp["products"] = products
	writeResponse(resp, w, http.StatusOK)
}

func HandlePostProducts(w http.ResponseWriter, r *http.Request) {
	// validate body contains valid product struct
	// TODO: this does not currently verify all required fields are present in json body
	resp := make(map[string]any)
	request_product := &models.NewProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&request_product)
	if err != nil {
		log.Print("Failed to create product")
		resp["error"] = err.Error()
		writeResponse(resp, w, http.StatusBadRequest)
		return
	}
	ok := request_product.ValidateProductRequest()
	if !ok {
		log.Print("Failed to create product")
		resp["error"] = err.Error()
		writeResponse(resp, w, http.StatusBadRequest)
		return
	}

	// try to do db insert
	new_product, err := product_repository.CreateProduct(request_product)
	if err != nil {
		log.Print("Failed to create product")
		resp["error"] = err.Error()
		writeResponse(resp, w, http.StatusBadGateway)
		return
	}
	// return new product to client
	resp["product"] = new_product
	writeResponse(resp, w, http.StatusCreated)
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
	writeResponse(resp, w, http.StatusOK)
}

func HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product, ok := ctx.Value("product").(*models.Product)
	if !ok {
		log.Print("Product not found in request context")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	err := product_repository.DeleteProduct(product)
	if err != nil {
		log.Print("Could not delete product")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func HandlePutProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	original_product, ok := ctx.Value("product").(*models.Product)
	if !ok {
		log.Print("Product not found in request context")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	request_product := &models.UpdateProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&request_product)
	if err != nil {
		log.Print("Failed to parse product update request from body")
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity) , http.StatusUnprocessableEntity)
		return
	}
	ok = request_product.ValidateProductRequest(original_product)
	if !ok {
		log.Print("Could not parse request body as product")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	updated_product, err := product_repository.UpdateProductById(*original_product, request_product)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	resp := make(map[string]any)
	resp["product"] = updated_product
	writeResponse(resp, w, http.StatusOK)
}

func writeResponse(resp map[string]any, w http.ResponseWriter, status_code int) {
	json_resp, err := json.Marshal(resp)
	if (err != nil) {
		log.Printf("Error occured in JSON marshal. Err %s\n", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	w.Write(json_resp)
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
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), "product", &product)
		next.ServeHTTP(w, r.WithContext(ctx))
  })
}