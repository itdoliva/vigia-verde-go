package handlerUser

import (
	"context"
	"encoding/json"
	"net/http"
	appUser "vigia-verde-go/internal/application/user"
	User "vigia-verde-go/internal/domain/user"
	repoUser "vigia-verde-go/internal/infrastructure/repository/user"
	web "vigia-verde-go/internal/infrastructure/utils"

	"cloud.google.com/go/firestore"
)

type UserService interface {
	CreateUser(ctx context.Context, dto appUser.RegisterReq) error
	GetByEmail(ctx context.Context, email string) (*User.User, error)
	GetById(ctx context.Context, id string) (*User.User, error)
	GetByPhone(ctx context.Context, phone string) (*User.User, error)
	Login(ctx context.Context, lr appUser.LoginReq) (string, error)
}

type UserHandler struct {
	service UserService
}

func NewHandler(s UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", h.Register)
	mux.HandleFunc("GET /getEmail/{email}", h.GetByEmail)
	mux.HandleFunc("GET /getId/{id}", h.GetById)
	mux.HandleFunc("GET /getPhone/{phone}", h.GetByPhone)
	mux.HandleFunc("POST /login", h.Login)
}
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto appUser.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(r.Context(), dto); err != nil {
		h.handleError(w, r, err)
		return
	}
	web.Respond(w, http.StatusOK, "Usuario Cadastrado")
}

func (h *UserHandler) GetByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	user, err := h.service.GetByEmail(r.Context(), email)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	web.Respond(w, http.StatusOK, user)
}
func (h *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.service.GetById(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	web.Respond(w, http.StatusOK, user)
}
func (h *UserHandler) GetByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.PathValue("phone")
	user, err := h.service.GetByPhone(r.Context(), phone)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	web.Respond(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto appUser.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	texto, err := h.service.Login(r.Context(), dto)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	web.Respond(w, http.StatusOK, texto)

}

func (h *UserHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	web.RespondError(w, r, err)
}

func SetupUser(db *firestore.Client) *UserHandler {
	repo := repoUser.NewRepository(db)
	service := appUser.NewService(repo)
	return NewHandler(service)
}
