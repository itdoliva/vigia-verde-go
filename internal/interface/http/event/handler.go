package handlerEvent

import (
	"context"
	"encoding/json"
	"net/http"
	appEvent "vigia-verde-go/internal/application/event"
	"vigia-verde-go/internal/domain/event"
	repoEvent "vigia-verde-go/internal/infrastructure/repository"
	web "vigia-verde-go/internal/infrastructure/utils"

	"cloud.google.com/go/firestore"
)

type Service interface {
	Create(ctx context.Context, input appEvent.CreateDTO) (string, error)
	ListAll(ctx context.Context, filter event.ListFilterParams) ([]event.ListEventResponse, int, error)
	GetByID(ctx context.Context, id string) (*event.Event, error)
}

type CreateRequest struct {
	Location struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	} `json:"location"`
	EventType   string `json:"event_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageSrc    string `json:"image_src"`
	AuthorID    string `json:"author_id"`
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
	mux.HandleFunc("GET /events/{id}", h.Get)
}

// Create cria um novo evento
// @Summary      Criar um evento
// @Description  Registra um novo evento a partir dos dados fornecidos no corpo da requisição.
// @Tags         Eventos
// @Accept       json
// @Produce      json
// @Param        request   body      appEvent.CreateDTO  true  "Dados para criação do evento"
// @Success      201       {object}  map[string]string   "Exemplo: {"id": "123"}"
// @Failure      400       {string}  string              "invalid json"
// @Failure      500       {object}  map[string]string   "Erro interno tratado"
// @Router       /events [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req appEvent.CreateDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	web.Respond(w, http.StatusCreated, map[string]string{"id": id})
}

// List retorna todos os eventos filtrados e paginados
// @Summary      Listar eventos
// @Description  Retorna uma lista de eventos com suporte a paginação, filtros por autor/tipo e busca geoespacial por coordenadas.
// @Tags         Eventos
// @Produce      json
// @Param        author_id   query     string  false  "Filtrar por ID do autor"
// @Param        event_type  query     string  false  "Filtrar por tipo de evento"
// @Param        page        query     int     false  "Número da página (padrão capturado por web.GetPagination)"
// @Param        limit       query     int     false  "Limite de itens por página"
// @Param        lat         query     number  false  "Latitude para busca geoespacial"
// @Param        lng         query     number  false  "Longitude para busca geoespacial"
// @Param        precision   query     string  false  "Precisão da busca (capturado por web.GetPrecision)"
// @Success      200         {array}   event.Event "Nota: O web.Respond envia a lista e o meta de paginação"
// @Failure      400         {string}  string      "Erro nos parâmetros enviados"
// @Failure      500         {object}  map[string]string
// @Router       /events [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	precision := web.GetPrecision(query)
	page, limit := web.GetPagination(query)
	latPtr, lngPtr, err := web.GetCoordinates(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := event.ListFilterParams{
		AuthorID:  query.Get("author_id"),
		EventType: query.Get("event_type"),
		Page:      page,
		Limit:     limit,
		Latitude:  latPtr,
		Longitude: lngPtr,
		Precision: precision,
	}

	events, total, err := h.service.ListAll(r.Context(), filter)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	meta := &web.PaginationMeta{TotalItems: total, CurrentPage: page}
	web.Respond(w, http.StatusOK, events, meta)
}

// Get retorna um evento específico pelo ID
// @Summary      Buscar um evento por ID
// @Description  Retorna os detalhes completos de um único evento baseado no parâmetro de URL {id}.
// @Tags         Eventos
// @Produce      json
// @Param        id        path      string  true  "ID do Evento"
// @Success      200       {object}  event.Event
// @Failure      404       {object}  map[string]string "Evento não encontrado"
// @Failure      500       {object}  map[string]string
// @Router       /events/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	event, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	web.Respond(w, http.StatusOK, event)
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	web.RespondError(w, r, err)
}

func SetupModule(db *firestore.Client) *Handler {
	repo := repoEvent.NewRepository(db)
	service := appEvent.NewService(repo)
	return NewHandler(service)
}
