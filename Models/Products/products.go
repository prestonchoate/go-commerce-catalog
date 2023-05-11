package products

import (
	"fmt"
	"log"

	models "github.com/prestonchoate/go-commerce-catalog/Models"
	services "github.com/prestonchoate/go-commerce-catalog/Services"
)

const TABLE_NAME = "products"

func GetAll() ([]models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v", TABLE_NAME)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Could not get data from %v", TABLE_NAME)
		log.Print(err.Error())
		log.Print(db.Stats())
		return nil, err
	}
	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := mapProductData(rows, &product)
		if err != nil {
			log.Print(err.Error())
			return products, fmt.Errorf("could not parse rows into products")
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProduct(product_id int) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v WHERE id = %v", TABLE_NAME, product_id)
	row := db.QueryRow(query)
	var product models.Product
	err := mapProductData(row, &product)
	if err != nil {
		log.Print(err.Error())
		return product, fmt.Errorf("failed to retrieve product with ID: %v", product_id)
	}
	return product, nil
}

func GetProductBySku(sku string) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v WHERE sku = \"%v\"", TABLE_NAME, sku)
	row := db.QueryRow(query)
	var product models.Product
	err := mapProductData(row, &product)
	if err != nil {
		log.Print(err.Error())
		return product, fmt.Errorf("failed to retrieve product with SKU: %v", sku)
	}
	return product, nil
}

func CreateProduct(input_product *models.NewProductRequest) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf(
		"INSERT INTO %v (sku, name, price, description) VALUES (%v, %v, %v, %v)",
		TABLE_NAME, 
		fmt.Sprintf("\"%v\"",input_product.SKU), 
		fmt.Sprintf("\"%v\"",input_product.Name), 
		fmt.Sprintf("\"%v\"",input_product.Price), 
		fmt.Sprintf("\"%v\"",input_product.Description),
	)
	result, err := db.Exec(query)
	if err != nil {
		blank_product := &models.Product{}
		log.Print(err.Error())
		return *blank_product, fmt.Errorf("could not create new product")
	}
	new_id, _ := result.LastInsertId()
	return GetProduct(int(new_id))
}

func DeleteProduct(product *models.Product) error {
	db := services.GetDB()
	query := fmt.Sprintf(
		"DELETE FROM %v WHERE ID = %v",
		TABLE_NAME,
		product.ID,
	)
	_, err := db.Exec(query)
	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("could not delete product with ID: %v", product.ID)
	}
	return nil
}

func UpdateProductById(original_product models.Product, input_product *models.UpdateProductRequest) (models.Product, error) {
	db := services.GetDB()
	query := fmt.Sprintf(
		"UPDATE %v SET sku = \"%v\", name = \"%v\", price = %v, description = \"%v\" WHERE id = %v",
		TABLE_NAME,
		input_product.SKU, 
		input_product.Name, 
		input_product.Price, 
		input_product.Description,
		original_product.ID)
	// TODO: Make this a transaction
	log.Printf("Attempting product update. Full query is: %s", query)
	_, err := db.Exec(query)
	if err != nil {
		log.Print(err.Error())
		return original_product, fmt.Errorf("could not update product")
	}
	return GetProduct(original_product.ID)
}

func mapProductData(row services.RowScanner, product *models.Product) ( error) {
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
		return fmt.Errorf("failed to map row data to product")
	}
	return nil
}