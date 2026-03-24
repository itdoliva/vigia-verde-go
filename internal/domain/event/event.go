package event

import (
	"context"
	"errors"
	"time"
)

type Repository interface {
	Create(ctx context.Context, ev *Event) (string, error)
	FindAll(ctx context.Context, filter ListFilterParams) ([]ListEventResponse, int, error)
	FindByID(ctx context.Context, id string) (*Event, error)
}

type Author struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Subscribed bool   `json:"subscribed"`
}

type EventType string

const (
	EventRemoval  EventType = "removal"
	EventPruning  EventType = "pruning"
	EventSeedling EventType = "seedling"
)

func (t EventType) Validate() error {
	switch t {
	case EventRemoval, EventPruning, EventSeedling:
		return nil
	default:
		return ErrInvalidEventType
	}
}

type GeoPoint struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type Event struct {
	ID           string    `json:"id"`
	Location     GeoPoint  `json:"location"`
	Geohash      string    `firestore:"geohash"`
	EventType    EventType `json:"event_type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Comments     []string  `json:"comments"`
	CommentCount int       `json:"comment_count"`
	ImageSrc     string    `json:"image_src"`
	Author       Author    `json:"author"`
	Upvotes      int       `json:"upvotes"`
	Downvotes    int       `json:"downvotes"`
	CreatedAt    time.Time `json:"created_at"`
}
type ListFilterParams struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Precision *int     `json:"precision"`
	AuthorID  string   `json:"author_id"`
	EventType string   `json:"event_type"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}
type ListEventResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       Author    `json:"author"`
	Location     GeoPoint  `json:"location"`
	EventType    EventType `json:"event_type"`
	CommentCount int       `json:"comment_count"`
	Upvotes      int       `json:"upvotes"`
	ImageSrc     string    `json:"image_src"`
	CreatedAt    time.Time `json:"created_at"`
}

func (f *ListFilterParams) Validate() error {
	if f.Latitude != nil || f.Longitude != nil {
		if f.Precision == nil {
			p := 6
			f.Precision = &p
		}
		if *f.Precision < 5 || *f.Precision > 8 {
			return errors.New("Invalid precision")
		}
	}
	return nil
}
