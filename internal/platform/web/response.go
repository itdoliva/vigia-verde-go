package web

type Response struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type PaginationMeta struct {
	TotalItems  int `json:"totalItems,omitempty"`
	CurrentPage int `json:"currentPage,omitempty"`
}
