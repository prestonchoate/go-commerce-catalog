package clients

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	models "github.com/prestonchoate/go-commerce-catalog/Models"
	services "github.com/prestonchoate/go-commerce-catalog/Services"
)

const TABLE_NAME = "clients"

func init() {
	db := services.GetDB()
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %v
		('id' int NOT NULL AUTO_INCREMENT,
		'client_id' varchar(255) NOT NULL,
		'client_secret' varchar(255) NOT NULL,
		'created_at' datetime DEFAULT CURRENT_TIMESTAMP,
		'last_login_time' datetime,
		PRIMARY KEY ('id'),
		UNIQUE KEY 'client_id' ('client_id'))`,
		TABLE_NAME)
	db.Exec(query)
}

func GetAll() ([]models.Client, error) {
	db := services.GetDB()
	query := fmt.Sprintf("SELECT * FROM %v", TABLE_NAME)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Could not get data from %v", TABLE_NAME)
		log.Print(err.Error())
		log.Print(db.Stats())
		return nil, err
	}
	var clients []models.Client
	for rows.Next() {
		client := &models.Client{}
		err := mapData(rows, *client)
		if err != nil {
			log.Print(err.Error())
			return clients, fmt.Errorf("could not parse rows into products")
		}
		clients = append(clients, *client)
	}
	return clients, nil
}

func mapData(r services.RowScanner, client models.Client) (error) {
	err := r.Scan(
		&client.ID,
		&client.Client_ID,
		&client.Client_Secret,
		&client.Created_At,
		&client.Last_Login_Time)
	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("failed to map client data to client object")
	}
	return nil
}

func hashClientSecret(secret string) (string, error) {
	// Convert string to byte slice
	var secretBytes = []byte(secret)

	// Hash password with Bcrypt's min cost
	hashedSecretBytes, err := bcrypt.GenerateFromPassword(secretBytes, bcrypt.MinCost)
	if err != nil {
		log.Print(err.Error())
		return "", fmt.Errorf("failed to hash client secret")
	}

	// Return hashed secret
	return string(hashedSecretBytes), nil
}

func isSecretCorrect(hashedSecret, currSecret string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedSecret), []byte(currSecret))
	return err == nil
}