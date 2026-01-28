package web

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data any             `json:"data"`
	Meta *PaginationMeta `json:"meta,omitempty"`
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
