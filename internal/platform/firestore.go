package platform

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

func initializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error initializing Firebase App: %w", err)
	}
	return app, nil
}

func createFirestoreClient(ctx context.Context, app *firebase.App) (*firestore.Client, error) {
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating Firestore client: %w", err)
	}
	return fsClient, nil
}

func SetupDatabaseConnection(ctx context.Context) (*firestore.Client, func(), error) {
	app, err := initializeFirebaseApp(ctx)
	if err != nil {
		return nil, nil, err
	}

	fsClient, err := createFirestoreClient(ctx, app)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to Firestore! ✅")

	cleanup := func() {
		if err := fsClient.Close(); err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		}
	}

	return fsClient, cleanup, nil
}
