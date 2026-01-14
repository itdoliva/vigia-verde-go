package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
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

func (r *EventRepository) FindAll(ctx context.Context, filter ListFilter) ([]Event, int, error) {
	collection := r.client.Collection("treeEvents")
	q := collection.Query

	if filter.AuthorID != "" {
		q = q.Where("authorId", "==", filter.AuthorID)
	}
	if filter.EventType != "" {
		q = q.Where("eventType", "==", filter.EventType)
	}

	aggRes, err := q.NewAggregationQuery().WithCount("all").Get(ctx)
	if err != nil {
		return nil, 0, err
	}

	raw, ok := aggRes["all"]
	if !ok {
		return nil, 0, errors.New(`missing aggregation key "all"`)
	}

	var totalCount int64
	switch v := raw.(type) {
	case int64:
		totalCount = v
	case *firestorepb.Value:
		totalCount = v.GetIntegerValue()
	default:
		return nil, 0, fmt.Errorf("unexpected type for count: %T", raw)
	}

	offset := (filter.Page - 1) * filter.Limit

	docs, err := q.Limit(filter.Limit).
		Offset(offset).
		Documents(ctx).GetAll()

	if err != nil {
		return nil, 0, err
	}

	events := make([]Event, 0, len(docs))
	for _, doc := range docs {
		var p persistenceModel
		if err := doc.DataTo(&p); err != nil {
			return nil, 0, err
		}

		events = append(events, Event{
			ID:        doc.Ref.ID,
			Location:  p.Location,
			EventType: EventType(p.EventType),
			Title:     p.Title,
			AuthorID:  p.AuthorID,
			Upvotes:   p.Upvotes,
			Downvotes: p.Downvotes,
			CreatedAt: p.CreatedAt,
		})
	}
	fmt.Printf("Filtros recebidos - Author: %s, Type: %s\n", filter.AuthorID, filter.EventType)
	return events, int(totalCount), nil
}
