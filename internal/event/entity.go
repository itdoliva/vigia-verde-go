package event

import (
	"errors"
	"net/http"
	"time"
	"vigia-verde-go/internal/platform/web"
)

var (
	ErrInvalidTitle     = web.Error{Err: errors.New("invalid title"), Status: http.StatusBadRequest}
	ErrInvalidEventType = web.Error{Err: errors.New("invalid eventType"), Status: http.StatusBadRequest}
	ErrInvalidGeoPoint  = web.Error{Err: errors.New("invalid location"), Status: http.StatusBadRequest}
	ErrInvalidPrecision = web.Error{Err: errors.New("invalid precision"), Status: http.StatusBadRequest}
	ErrNotFound         = web.Error{Err: errors.New("event not found"), Status: http.StatusNotFound}
)

type GeoPoint struct {
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
}

func (p GeoPoint) Validate() error {
	if p.Latitude < -90 || p.Latitude > 90 {
		return ErrInvalidGeoPoint
	}
	if p.Longitude < -180 || p.Longitude > 180 {
		return ErrInvalidGeoPoint
	}
	return nil
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

type CreateInput struct {
	Location    GeoPoint
	EventType   EventType
	Title       string
	Description string
	ImageSrc    string
	AuthorID    string
}

type Author struct {
	ID         string `json:"id" firestore:"id"`
	Name       string `json:"name" firestore:"name"`
	Subscribed bool   `json:"subscribed" firestore:"subscribed"`
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

func New(in CreateInput) (*Event, error) {
	if in.Title == "" {
		return nil, ErrInvalidTitle
	}
	if err := in.Location.Validate(); err != nil {
		return nil, err
	}
	if err := in.EventType.Validate(); err != nil {
		return nil, err
	}
	if in.AuthorID == "" {
		return nil, errors.New("author id is required")
	}

	return &Event{
		Location:    in.Location,
		EventType:   in.EventType,
		Title:       in.Title,
		Description: in.Description,
		ImageSrc:    in.ImageSrc,
		Author: Author{
			ID:         in.AuthorID,
			Name:       "test",
			Subscribed: true,
		},
		Upvotes:      0,
		Downvotes:    0,
		CommentCount: 0,
		Comments:     []string{},
	}, nil
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
