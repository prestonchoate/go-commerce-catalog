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
		return product, fmt.Errorf("Failed to retrieve product with SKU: %v", sku)
	}
	return product, nil
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