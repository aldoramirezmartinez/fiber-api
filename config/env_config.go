package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %s\n", err)
		return err
	}
	return nil
}

func GetPort() string {
	return os.Getenv("PORT")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}

func GetMongoURI() (string, error) {
	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		err := fmt.Errorf("MongoDB URI not found in .env file")
		fmt.Println(err)
		return "", err
	}

	return mongoURI, nil
}
