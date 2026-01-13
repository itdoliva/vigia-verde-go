package event

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type EventRepository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *EventRepository {
	return &EventRepository{client: client}
}

func (r *EventRepository) Create(ctx context.Context, ev *Event) (string, error) {
	collection := r.client.Collection("treeEvents")
	docRef := collection.NewDoc()

	data := persistenceModel{
		Location:  ev.Location,
		EventType: string(ev.EventType),
		Title:     ev.Title,
		AuthorID:  ev.AuthorID,
		Upvotes:   ev.Upvotes,
		Downvotes: ev.Downvotes,
	}

	if _, err := docRef.Set(ctx, data); err != nil {
		return "", err
	}

	ev.ID = docRef.ID
	return docRef.ID, nil
}

type persistenceModel struct {
	Location  GeoPoint  `firestore:"location"`
	EventType string    `firestore:"eventType"`
	Title     string    `firestore:"title"`
	AuthorID  string    `firestore:"authorId"`
	Upvotes   int       `firestore:"upvotes"`
	Downvotes int       `firestore:"downvotes"`
	CreatedAt time.Time `firestore:"createdAt,serverTimestamp"`
}
