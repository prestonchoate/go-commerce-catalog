package models

import (
	"log"
	"time"
)

type Product struct {
	ID int `json:"id"`
	SKU string `json:"sku"`
	Name string `json:"name"`
	Price float32 `json:"price"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewProductRequest struct {
	SKU string `json:"sku"`
	Name string `json:"name"`
	Price float32 `json:"price"`
	Description string `json:"description"`
}

func (product_request *NewProductRequest) ValidateProductRequest() bool {
	if (
		len(product_request.SKU) == 0 ||
		len(product_request.Name) == 0 ||
		len(product_request.Description) == 0 ||
		product_request.Price <= 0.0) {
			return false
		}
		return true
}
type UpdateProductRequest struct {
		SKU string `json:"sku"`
		Name string `json:"name"`
		Price float32 `json:"price"`
		Description string `json:"description"`
}

func (product_request *UpdateProductRequest) ValidateProductRequest(original_product *Product) bool {
	if len(product_request.SKU) == 0 {
		log.Print("Setting original data on product update request for SKU")
		product_request.SKU = original_product.SKU
	}
	
	if len(product_request.Name) == 0 {
		log.Print("Setting original data on product update request for Name")
		product_request.Name = original_product.Name
	}

	if len(product_request.Description) == 0 {
		log.Print("Setting original data on product update request for Description")
		product_request.Description = original_product.Description
	}

	if product_request.Price <= 0.0 {
		log.Print("Setting original data on product update request for Price")
		product_request.Price = original_product.Price
	}

	return true;
}