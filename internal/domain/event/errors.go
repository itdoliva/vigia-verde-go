package event

import (
	"errors"
	"net/http"
	web "vigia-verde-go/internal/infrastructure/utils"
)

var (
	ErrInvalidTitle     = web.Error{Err: errors.New("invalid title"), Status: http.StatusBadRequest}
	ErrInvalidEventType = web.Error{Err: errors.New("invalid eventType"), Status: http.StatusBadRequest}
	ErrInvalidGeoPoint  = web.Error{Err: errors.New("invalid location"), Status: http.StatusBadRequest}
	ErrInvalidPrecision = web.Error{Err: errors.New("invalid precision"), Status: http.StatusBadRequest}
	ErrInvalidID        = web.Error{Err: errors.New("author id is required"), Status: http.StatusBadRequest}
	ErrNotFound         = web.Error{Err: errors.New("event not found"), Status: http.StatusNotFound}
)
