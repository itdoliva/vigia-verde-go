package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Response struct {
	Data any             `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginationMeta struct {
	TotalItems  int `json:"total_items,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
}

func Respond(w http.ResponseWriter, status int, data any, meta ...*PaginationMeta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := Response{
		Data: data,
	}

	if len(meta) > 0 && meta[0] != nil {
		res.Meta = meta[0]

	}

	json.NewEncoder(w).Encode(res)

}

// No package web

// Error representa um erro que sabe seu status HTTP e sua mensagem de resposta.
type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

// RespondError agora centraliza a lógica de log e resposta
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	// Verifica se é um erro "conhecido" do nosso sistema
	var webErr Error
	if errors.As(err, &webErr) {
		status = webErr.Status
		message = webErr.Error()
	}

	// Log centralizado (ajuda a manter o console limpo e padronizado)
	log.Printf("[ERROR] %s %s: %v", r.Method, r.URL.Path, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
