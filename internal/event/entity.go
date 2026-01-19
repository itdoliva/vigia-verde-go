package event

import (
	"errors"
	"math"
	"time"
)

var (
	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidEventType = errors.New("invalid eventType")
	ErrInvalidGeoPoint  = errors.New("invalid location")
)

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	EventCut      EventType = "cut"
	EventPruning  EventType = "pruning"
	EventPlanting EventType = "planting"
)

func (t EventType) Validate() error {
	switch t {
	case EventCut, EventPruning, EventPlanting:
		return nil
	default:
		return ErrInvalidEventType
	}
}

type CreateInput struct {
	Location  GeoPoint
	EventType EventType
	Title     string
	AuthorID  string
}

type Event struct {
	ID        string    `json:"id"`
	Location  GeoPoint  `json:"location"`
	EventType EventType `json:"event_type"`
	Title     string    `json:"title"`
	AuthorID  string    `json:"author_id"`
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
	CreatedAt time.Time `json:"created_at"`
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

	return &Event{
		Location:  in.Location,
		EventType: in.EventType,
		Title:     in.Title,
		AuthorID:  in.AuthorID,
		Upvotes:   0,
		Downvotes: 0,
	}, nil
}

type ListFilter struct {
	Latitude  *float64
	Longitude *float64
	Radius    float64
	AuthorID  string
	EventType string
	Page      int
	Limit     int
}

// Formula de Haversine
func (p GeoPoint) Distancia(other GeoPoint) float64 {
	const R = 6371000 // Raio da Terra em metros
	phi1 := p.Latitude * math.Pi / 180
	phi2 := other.Latitude * math.Pi / 180
	dphi := (other.Latitude - p.Latitude) * math.Pi / 180
	dlng := (other.Longitude - p.Longitude) * math.Pi / 180

	a := math.Sin(dphi/2)*math.Sin(dphi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
