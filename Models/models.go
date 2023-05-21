package models

import (
	"log"
	"time"
)

type Client struct {
	ID int `json:"id"`
	Client_ID string `json:"clientId"`
	Client_Secret string `json:"clientSecret"`
	Created_At time.Time `json:"createdAt"`
	Last_Login_Time time.Time `json:"lastLoginTime"`
}

type ClientRequest struct {
	Client_ID string `json:"clientId"`
	Client_Secret string `json:"clientSecret"`
}

type Product struct {
	ID int `json:"id"`
	SKU string `json:"sku"`
	Name string `json:"name"`
	Price float32 `json:"price"`
	Description string `json:"description"`
	Created_At time.Time `json:"createdAt"`
	Updated_At time.Time `json:"updatedAt"`
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