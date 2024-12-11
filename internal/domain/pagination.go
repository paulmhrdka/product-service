package domain

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PageSize    int   `json:"page_size"`
}
