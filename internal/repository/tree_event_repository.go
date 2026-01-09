package repository

import (
	"context"
	"vigia-verde-go/internal/core"

	"cloud.google.com/go/firestore"
)

type TreeEventRepository struct {
	client *firestore.Client
}

func NewTreeEventRepository(client *firestore.Client) *TreeEventRepository {
	return &TreeEventRepository{client: client}
}

func (r *TreeEventRepository) Create(ctx context.Context, input core.TreeEvent) (string, error) {
	collection := r.client.Collection("treeEvents")
	ref := collection.NewDoc()

	doc := map[string]interface{}{
		"location": map[string]float64{
			"latitude":  input.Location.Latitude,
			"longitude": input.Location.Longitude,
		},
		"eventType": string(input.EventType),
		"title":     input.Title,
		"authorId":  input.AuthorID,
		"upvotes":   input.Upvotes,
		"downvotes": input.Downvotes,
		"createdAt": firestore.ServerTimestamp,
	}

	_, err := ref.Set(ctx, doc)
	if err != nil {
		return "", err
	}

	return ref.ID, nil
}
