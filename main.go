package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type S3bucket struct {
	BucketName string
	S3Region   string
	// uploadDir  string
	S3Client *s3.Client
}

type application struct {
	S3bucket
}

func main() {
	if os.Getenv("ENV") != "Production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

	}

	app := &application{}

	app.BucketName = os.Getenv("BUCKETNAME")
	app.S3Region = os.Getenv("S3REGION")

	if app.BucketName == "" || app.S3Region == "" {
		log.Println("Usage: go run main.go -bucket=<s3_bucket> -regiion=<aws_region> -uploadDir=<aws_store_directory>")
		os.Exit(1)
	}

	//load AWS config
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(app.S3Region))
	if err != nil {
		panic("failed to load AWS config")
	}

	app.S3Client = s3.NewFromConfig(cfg)

	mux := mux.NewRouter()

	mux.HandleFunc("/signedurl/{filename}", app.generatePresignedUrl).Methods(http.MethodGet)

	log.Println("server started")
	http.ListenAndServe(":8080", mux)

}
