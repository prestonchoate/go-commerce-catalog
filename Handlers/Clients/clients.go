package clients

import (
	"encoding/json"
	"log"
	"net/http"

	client_repository "github.com/prestonchoate/go-commerce-catalog/Models/Clients"
)


func HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]any)
	clients, err := client_repository.GetAll()
	if err != nil {
		log.Print("Failed to retrieve clients")
		resp["error"] = err.Error()
		writeResponse(resp, w, http.StatusBadGateway)
		return
	}
	resp["clients"] = clients
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