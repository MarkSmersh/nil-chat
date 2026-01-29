package utils

import (
	"log"
	"os"
)

func GetJWTSecret() string {
	secret := os.Getenv("SECRET")

	if secret == "" {
		log.Fatal("The SECRET enviroment variable is empty")
	}

	return secret
}

func GetHashSecret() string {
	secret := os.Getenv("SECRET")

	if secret == "" {
		log.Fatal("The SECRET enviroment variable is empty")
	}

	if len(secret) != 32 {
		log.Fatal("The SECRET enviroment variable needs to have a length of 32 symbols")
	}

	return secret
}

func GetDBUrl() string {
	dburl := os.Getenv("DB_URL")

	if dburl == "" {
		log.Fatal("The DB_URL environment is missing")
	}

	return dburl
}
