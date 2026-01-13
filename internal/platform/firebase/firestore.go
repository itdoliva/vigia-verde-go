package firebase

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func NewFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	jsonPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if jsonPath == "" {
		return nil, fmt.Errorf("environment variable GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, jsonPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return client, nil
}
