package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	var webErr Error
	if errors.As(err, &webErr) {
		status = webErr.Status
		message = webErr.Error()
	}

	log.Printf("[ERROR] %s %s: %v", r.Method, r.URL.Path, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
