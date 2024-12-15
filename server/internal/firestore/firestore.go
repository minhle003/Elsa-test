package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func NewClient() (*firestore.Client, error) {
	ctx := context.Background()

	// Ensure the service account key file path is set in the environment variable
	serviceAccountKeyPath := "/Users/minh.lequang/Elsa-test/server/dbKey.json"
	if serviceAccountKeyPath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
	}

	// Create Firestore client using the service account key file
	client, err := firestore.NewClient(ctx, "englishquizappdb", option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return nil, err
	}
	return client, nil
}
