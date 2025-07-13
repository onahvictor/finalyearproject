package main

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
)

func (app *application) generatePresignedUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Extract the parameters from the URL
	filename := vars["filename"]
	s3Client := s3.NewPresignClient(app.S3Client)

	req, err := s3Client.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(app.BucketName),
		Key:         aws.String(filepath.Join("audio", filename)),
		ContentType: aws.String("audio/wav"), // Add this
	})

	if err != nil {
		http.Error(w, "Failed to generate presigned URL", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(req.URL))
}
