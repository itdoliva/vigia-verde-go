package treeeventhandler

import (
	"encoding/json"
	"log"
	"net/http"

	treeevent "vigia-verde-go/internal/core"
	"vigia-verde-go/internal/service"
)

type TreeEventHandler struct {
	service service.TreeEventService
}

func NewTreeEventHandler(s service.TreeEventService) *TreeEventHandler {
	return &TreeEventHandler{
		service: s,
	}
}

func (h *TreeEventHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tree-events", h.handleCreateTreeEvent)
}

func (h *TreeEventHandler) handleCreateTreeEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var input treeevent.CreateTreeEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON body",
		})
		return
	}
	defer r.Body.Close()

	id, err := h.service.CreateTreeEvent(r.Context(), input)
	if err != nil {
		log.Printf("erro ao criar treeEvent: %v", err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to create treeEvent",
		})
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	})
}
