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

func (r *EventRepository) FindAll(ctx context.Context) ([]Event, error) {
	collection := r.client.Collection("treeEvents")
	docs, err := collection.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	events := make([]Event, 0, len(docs))
	for _, doc := range docs {
		var p persistenceModel
		if err := doc.DataTo(&p); err != nil {
			return nil, err
		}

		events = append(events, Event{
			ID:        doc.Ref.ID,
			Location:  p.Location,
			EventType: EventType(p.EventType),
			Title:     p.Title,
			AuthorID:  p.AuthorID,
			Upvotes:   p.Upvotes,
			Downvotes: p.Downvotes,
			CreateAt:  p.CreatedAt,
		})
	}
	return events, nil

}
