package event

import (
	"time"
)

type Author struct {
	ID         string `json:"id" firestore:"id"`
	Name       string `json:"name" firestore:"name"`
	Subscribed bool   `json:"subscribed" firestore:"subscribed"`
}

type EventType string

const (
	EventRemoval  EventType = "removal"
	EventPruning  EventType = "pruning"
	EventSeedling EventType = "seedling"
)

type GeoPoint struct {
	Latitude  *float64 `json:"latitude" firestore:"latitude"`
	Longitude *float64 `json:"longitude" firestore:"longitude"`
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

func (t EventType) Validate() error {
	switch t {
	case EventRemoval, EventPruning, EventSeedling:
		return nil
	default:
		return ErrInvalidEventType
	}
}

type EventDto struct {
	Location    GeoPoint
	EventType   EventType
	Title       string
	Description string
	ImageSrc    string
	AuthorID    string
}

type ListFilterParams struct {
	Latitude  *float64
	Longitude *float64
	Precision *int
	AuthorID  string
	EventType string
	Page      int
	Limit     int
}

func (f *ListFilterParams) Validate() error {
	if f.Latitude != nil || f.Longitude != nil {
		if f.Precision == nil {
			p := 6
			f.Precision = &p
		}
		if *f.Precision < 5 || *f.Precision > 8 {
			return ErrInvalidPrecision
		}
	}
	return nil
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
