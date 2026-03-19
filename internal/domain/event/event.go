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
