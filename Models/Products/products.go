package products

import (
	"fmt"
	"log"

	models "github.com/prestonchoate/go-commerce-catalog/Models"
	services "github.com/prestonchoate/go-commerce-catalog/Services"
)

var table_name = "products"

func GetAll() ([]models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v", table_name)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Could not get data from %v", table_name)
		log.Print(err.Error())
		log.Print(db.Stats())
		return nil, err
	}
	var products []models.Product
	for rows.Next() {
		var product models.Product
		product, err := mapProductData(rows, &product)
		if err != nil {
			log.Print(err.Error())
			return products, fmt.Errorf("Could not parse rows into products")
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProduct(product_id int) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v WHERE id = %v", table_name, product_id)
	row := db.QueryRow(query)
	var product models.Product
	product, err := mapProductData(row, &product)
	if err != nil {
		log.Print(err.Error())
		return product, fmt.Errorf("failed to retrieve product with ID: %v", product_id)
	}
	return product, nil
}

func GetProductBySku(sku string) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v WHERE sku = %v", table_name, sku)
	row := db.QueryRow(query)
	var product models.Product
	product, err := mapProductData(row, &product)
	if err != nil {
		log.Print(err.Error())
		return product, fmt.Errorf("failed to retrieve product with SKU: %v", sku)
	}
	return product, nil
}

func UpdateProductById(original_product models.Product, input_product models.Product) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf(
		"INSERT INTO %v (sku, name, price, description) VALUES (%v, %v, %v, %v) WHERE id = %v",
		table_name, 
		input_product.SKU, 
		input_product.Name, 
		input_product.Price, 
		input_product.Description,
		original_product.ID)
	// TODO: Make this a transaction
	// TODO: Check result to make sure the updated ID is correct
	_, err := db.Exec(query)
	if err != nil {
		log.Print(err.Error())
		return input_product, fmt.Errorf("could not update product")
	}
	return input_product, nil
}

func mapProductData(row services.RowScanner, product *models.Product) (models.Product, error) {
	err := row.Scan(
		&product.ID,
		&product.SKU,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt)
	if err != nil {
		log.Print(err.Error())
		return *product, fmt.Errorf("Failed to map row data to product")
	}
	return *product, nil
}