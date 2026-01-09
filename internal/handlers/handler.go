package handlers

import (
	"encoding/json"
	"net/http"

	"vigia-verde-go/internal/core"
	service "vigia-verde-go/internal/service"
)

type CreateTreeEventRequest struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	EventType string `json:"eventType"`
	Title     string `json:"title"`
	AuthorID  string `json:"authorId"`
}

type TreeEventHandler struct {
	service *service.TreeEventService
}

func NewTreeEventHandler(s *service.TreeEventService) *TreeEventHandler {
	return &TreeEventHandler{service: s}
}

func (h *TreeEventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tree-events", h.handleTreeEvents)
}

func (h *TreeEventHandler) handleTreeEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Create(w, r)
	default:
		w.Header().Set("Allow", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TreeEventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTreeEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	input := core.CreateInput{
		Location: core.GeoPoint{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		EventType: core.EventType(req.EventType),
		Title:     req.Title,
		AuthorID:  req.AuthorID,
	}

	id, err := h.service.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}
