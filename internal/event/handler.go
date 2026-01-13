package event

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (string, error)
}

type CreateRequest struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	EventType string `json:"eventType"`
	Title     string `json:"title"`
	AuthorID  string `json:"authorId"`
}

func (req CreateRequest) toInput() CreateInput {
	return CreateInput{
		Location: GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		EventType: EventType(req.EventType),
		Title:     req.Title,
		AuthorID:  req.AuthorID,
	}
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /events", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), req.toInput())
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	log.Print(err)

	switch {
	case errors.Is(err, ErrInvalidTitle),
		errors.Is(err, ErrInvalidEventType),
		errors.Is(err, ErrInvalidGeoPoint):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func SetupModule(db *firestore.Client) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return NewHandler(service)
}
