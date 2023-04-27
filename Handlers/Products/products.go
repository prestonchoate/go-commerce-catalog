package products

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	product_repository "github.com/prestonchoate/go-commerce-catalog/Models/Products"
)

func HandleProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling %v to %v", r.Method, r.RequestURI)
	if r.Method == http.MethodGet {
		handleGetAllProducts(w, r)
	}
}

func handleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	products, err := product_repository.GetAll()
	if err != nil {
		log.Fatal("Failed to retrieve products")
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

func generateResponse(resp map[string]any) ([]byte, error){
	json_resp, err := json.Marshal(resp)
	if (err != nil) {
		log.Fatalf("Error occured in JSON marshal. Err: %s\n", err)
		return nil, errors.New("could not marshal JSON response")
	}
	return json_resp, nil
}