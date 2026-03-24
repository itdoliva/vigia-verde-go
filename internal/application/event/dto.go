package appEvent

import (
	"vigia-verde-go/internal/domain/event"

	"github.com/go-playground/validator/v10"
)

type CreateDTO struct {
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
	EventType   string  `json:"event_type" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	ImageSrc    string  `json:"image_src"`
}
type UpdateDTO struct {
	Latitude    *float64 `json:"latitude" validate:"latitude"`
	Longitude   *float64 `json:"longitude" validate:"longitude"`
	EventType   *string  `json:"event_type"`
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	ImageSrc    *string  `json:"image_src"`
}

func New(in CreateDTO) (*event.Event, error) {
	if err := validator.New().Struct(in); err != nil {
		return nil, err
	}
	eventType := event.EventType(in.EventType)
	geopoint := &event.GeoPoint{
		Latitude:  &in.Latitude,
		Longitude: &in.Longitude,
	}
	if in.Title == "" {
		return nil, event.ErrInvalidTitle
	}
	if err := eventType.Validate(); err != nil {
		return nil, err
	}

	return &event.Event{
		Location:    *geopoint,
		EventType:   event.EventType(in.EventType),
		Title:       in.Title,
		Description: in.Description,
		ImageSrc:    in.ImageSrc,
		Author: event.Author{
			ID:         "",
			Name:       "",
			Subscribed: true,
		},
		Upvotes:      0,
		Downvotes:    0,
		CommentCount: 0,
		Comments:     []string{},
	}, nil
}
