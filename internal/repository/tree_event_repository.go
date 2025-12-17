package repository

import (
	"context"
	treeevent "vigia-verde-go/internal/core"

	"cloud.google.com/go/firestore"
)

type TreeEventRepository interface {
	Create(ctx context.Context, input treeevent.CreateTreeEventInput) (string, error)
}

type treeEventRepository struct {
	client *firestore.Client
}

func NewTreeEventRepository(client *firestore.Client) TreeEventRepository {
	return &treeEventRepository{client: client}
}

func (r *treeEventRepository) Create(ctx context.Context, input treeevent.CreateTreeEventInput) (string, error) {
	col := r.client.Collection("treeEvents")
	ref := col.NewDoc()
	doc := map[string]interface{}{
		"location": map[string]float64{
			"latitude":  input.Location.Latitude,
			"longitude": input.Location.Longitude,
		},
		"eventType": input.EventType,
		"title":     input.Title,
		"authorId":  input.AuthorID,
		"upvotes":   0,
		"downvotes": 0,
		"createdAt": firestore.ServerTimestamp,
	}

	_, err := ref.Set(ctx, doc)
	if err != nil {
		return "", err
	}

	return ref.ID, nil
}
