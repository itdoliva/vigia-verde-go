package web

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type PaginationMeta struct {
	TotalItems  int `json:"total_items,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
}
