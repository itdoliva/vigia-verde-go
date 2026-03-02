package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/mmcloughlin/geohash"
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
	hash := geohash.Encode(ev.Location.Latitude, ev.Location.Longitude)
	tokens := []string{
		hash[:5],
		hash[:6],
		hash[:7],
		hash[:8],
	}

	data := persistenceModel{
		Location:      ev.Location,
		Geohash:       hash,
		GeohashTokens: tokens,
		EventType:     string(ev.EventType),
		Title:         ev.Title,
		AuthorID:      ev.AuthorID,
		Upvotes:       ev.Upvotes,
		Downvotes:     ev.Downvotes,
	}

	if _, err := docRef.Set(ctx, data); err != nil {
		return "", err
	}

	ev.ID = docRef.ID
	return docRef.ID, nil
}

// Falta adicionar URL da imagem, comentarios e numero de comentarios
type persistenceModel struct {
	Location      GeoPoint  `firestore:"location"`
	Geohash       string    `firestore:"geohash"`
	GeohashTokens []string  `firestore:"geohashTokens"`
	EventType     string    `firestore:"eventType"`
	Title         string    `firestore:"title"`
	AuthorID      string    `firestore:"authorId"`
	Upvotes       int       `firestore:"upvotes"`
	Downvotes     int       `firestore:"downvotes"`
	CreatedAt     time.Time `firestore:"createdAt,serverTimestamp"`
}

func (r *EventRepository) FindAll(ctx context.Context, filter ListFilterParams) ([]ListEventResponse, int, error) {
	collection := r.client.Collection("treeEvents")
	q := collection.Query

	if filter.AuthorID != "" {
		q = q.Where("authorId", "==", filter.AuthorID)
	}
	if filter.EventType != "" {
		q = q.Where("eventType", "==", filter.EventType)
	}

	if filter.Latitude != nil && filter.Longitude != nil {
		fullHash := geohash.Encode(*filter.Latitude, *filter.Longitude)

		precision := 6
		if filter.Precision != nil {
			precision = *filter.Precision
		}

		centerHash := fullHash[:precision]
		prefixes := append(geohash.Neighbors(centerHash), centerHash)
		q = q.Where("geohashTokens", "array-contains-any", prefixes)
	}

	aggregation, err := q.NewAggregationQuery().WithCount("total_count").Get(ctx)
	if err != nil {
		return nil, 0, err
	}

	countValue, ok := aggregation["total_count"]
	if !ok {
		return nil, 0, errors.New(`missing aggregation key "total_count"`)
	}

	var totalCount int64
	switch v := countValue.(type) {
	case int64:
		totalCount = v
	case *firestorepb.Value:
		totalCount = v.GetIntegerValue()
	default:
		return nil, 0, fmt.Errorf("unexpected type for count: %T", countValue)
	}

	offset := (filter.Page - 1) * filter.Limit

	docs, err := q.Limit(filter.Limit).Offset(offset).Documents(ctx).GetAll()

	if err != nil {
		return nil, 0, err
	}

	events := make([]ListEventResponse, 0, len(docs))
	for _, doc := range docs {
		var p persistenceModel
		if err := doc.DataTo(&p); err != nil {
			return nil, 0, err
		}

		events = append(events, ListEventResponse{
			ID:        doc.Ref.ID,
			Location:  p.Location,
			EventType: EventType(p.EventType),
			CreatedAt: p.CreatedAt,
		})
	}
	return events, int(totalCount), nil
}

func (r *EventRepository) FindByID(ctx context.Context, id string) (*Event, error) {
	doc, err := r.client.Collection("treeEvents").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var p persistenceModel
	if err := doc.DataTo(&p); err != nil {
		return nil, err
	}

	return &Event{
		ID:        doc.Ref.ID,
		Location:  p.Location,
		EventType: EventType(p.EventType),
		Title:     p.Title,
		AuthorID:  p.AuthorID,
		Upvotes:   p.Upvotes,
		Downvotes: p.Downvotes,
		CreatedAt: p.CreatedAt,
	}, nil
}
