package event

import (
	"errors"
	"time"
)

var (
	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidEventType = errors.New("invalid eventType")
	ErrInvalidGeoPoint  = errors.New("invalid location")
)

type GeoPoint struct {
	Latitude  float64
	Longitude float64
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
	ID        string
	Location  GeoPoint
	EventType EventType
	Title     string
	AuthorID  string
	Upvotes   int
	Downvotes int
	CreateAt  time.Time
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
