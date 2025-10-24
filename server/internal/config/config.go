package config

import (
	"log"
	"os"
)

type config struct {
	MongoDB_url    string
	PORT           string
	JWT_SECRET_KEY string
}

func LoadConfig() config {

	db := os.Getenv("DB_CONNECTION_URL")

	if db == "" {
		log.Fatalf("DB_CONNECTION_URL not set")
	}

	port := os.Getenv("PORT")

	if port == "" {

		log.Fatalf("PORT not set")

	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	if jwtSecretKey == "" {
		log.Fatalf(" JWT_SECRET_KEY not set")
	}

	return config{
		MongoDB_url:    db,
		PORT:           port,
		JWT_SECRET_KEY: jwtSecretKey,
	}

}
