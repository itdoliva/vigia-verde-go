package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"vigia-verde-go/internal/domain/event"
)

type Response[T any] struct {
	Data T               `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginationMeta struct {
	TotalItems  int `json:"total_items,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
}

func Respond[T any](w http.ResponseWriter, status int, data T, meta ...*PaginationMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := Response[T]{
		Data: data,
	}

	if len(meta) > 0 && meta[0] != nil {
		res.Meta = meta[0]

	}

	json.NewEncoder(w).Encode(res)

}

// domainErrors mapeia erros de domínio para HTTP status + mensagem pública.
var domainErrors = map[error]struct {
	Status  int
	Message string
}{
	event.ErrNotFound:         {http.StatusNotFound, "event not found"},
	event.ErrInvalidTitle:     {http.StatusBadRequest, "invalid title"},
	event.ErrInvalidEventType: {http.StatusBadRequest, "invalid event type"},
	event.ErrInvalidGeoPoint:  {http.StatusBadRequest, "invalid location"},
	event.ErrInvalidPrecision: {http.StatusBadRequest, "invalid precision"},
	event.ErrInvalidID:        {http.StatusBadRequest, "author id is required"},
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	for domainErr, meta := range domainErrors {
		if errors.Is(err, domainErr) {
			Respond(w, meta.Status, ErrorResponse{Error: meta.Message})
			return
		}
	}

	log.Printf("[ERROR] %s %s: %v", r.Method, r.URL.Path, err)
	Respond(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}
