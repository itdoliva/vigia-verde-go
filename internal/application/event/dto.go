package appEvent

import (
	"errors"
	"time"
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

type ListFilterParams struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Precision *int     `json:"precision"`
	AuthorID  string   `json:"author_id"`
	EventType string   `json:"event_type"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
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

type ListEventResponse struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Author       event.Author    `json:"author"`
	Location     event.GeoPoint  `json:"location"`
	EventType    event.EventType `json:"event_type"`
	CommentCount int             `json:"comment_count"`
	Upvotes      int             `json:"upvotes"`
	ImageSrc     string          `json:"image_src"`
	CreatedAt    time.Time       `json:"created_at"`
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
