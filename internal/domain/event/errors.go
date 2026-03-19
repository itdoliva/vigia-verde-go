package event

import (
	"errors"
)

var (
	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidEventType = errors.New("invalid eventType")
	ErrInvalidGeoPoint  = errors.New("invalid location")
	ErrInvalidPrecision = errors.New("invalid precision")
	ErrInvalidID        = errors.New("author id is required")
	ErrNotFound         = errors.New("event not found")
)
