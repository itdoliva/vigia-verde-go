package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vigia-verde-go/internal/platform/web"

	"cloud.google.com/go/firestore"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	ListAll(ctx context.Context, filter ListFilter) ([]Event, int, error)
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
	mux.HandleFunc("GET /events", h.List)
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
	_ = json.NewEncoder(w).Encode(web.Response{
		Data: map[string]string{"id": id},
	})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 30
	}
	if limit > 50 {
		limit = 50
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	var latPtr, lngPtr *float64
	latStr := query.Get("lat")
	lngStr := query.Get("lng")
	if latStr != "" && lngStr != "" {
		l1, _ := strconv.ParseFloat(latStr, 64)
		l2, _ := strconv.ParseFloat(lngStr, 64)
		latPtr = &l1
		lngPtr = &l2
	}

	radius := 100.0
	if rStr := query.Get("radius"); rStr != "" {
		if rVal, err := strconv.ParseFloat(rStr, 64); err == nil {
			radius = rVal
		}
	}

	filter := ListFilter{
		AuthorID:  query.Get("authorId"),
		EventType: query.Get("eventType"),
		Page:      page,
		Limit:     limit,
		Latitude:  latPtr,
		Longitude: lngPtr,
		Radius:    radius,
	}

	events, total, err := h.service.ListAll(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(web.Response{
		Data: events,
		Meta: web.PaginationMeta{
			TotalItems:  total,
			CurrentPage: page,
		},
	})

	fmt.Println("URL Completa recebida:", r.URL.String())
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	log.Printf("[ERROR] %v", err)
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
