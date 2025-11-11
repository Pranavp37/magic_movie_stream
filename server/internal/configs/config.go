package configs

import (
	"log"
	"os"
)

type configs struct {
	MongoDB_url           string
	PORT                  string
	JWT_SECRET_KEY        string
	BucketName            string
	S3_Region             string
	AWS_SECRET_ACCESS_KEY string
	AWS_ACCESS_KEY_ID     string
}

func LoadConfig() configs {

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

	bucketname := os.Getenv("BUCKET_NAME")
	if bucketname == "" {
		log.Fatalf("BUCKET_NAME not set")
	}

	s3_region := os.Getenv("S3_REGION")
	if s3_region == "" {
		log.Fatalf("S3_REGION not set")
	}

	aws_access_id := os.Getenv("AWS_ACCESS_KEY_ID")
	if aws_access_id == "" {
		log.Fatalf("AWS_ACCESS_KEY_ID not set")
	}

	aws_secret_access_key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if aws_secret_access_key == "" {
		log.Fatalf("AWS_SECRET_ACCESS_KEY not set")

	}

	return configs{
		MongoDB_url:           db,
		PORT:                  port,
		JWT_SECRET_KEY:        jwtSecretKey,
		BucketName:            bucketname,
		S3_Region:             s3_region,
		AWS_SECRET_ACCESS_KEY: aws_secret_access_key,
		AWS_ACCESS_KEY_ID:     aws_access_id,
	}

}
