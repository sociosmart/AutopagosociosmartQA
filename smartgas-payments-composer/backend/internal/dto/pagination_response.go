package dto

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Data       any   `json:"data" format:"array"`
}
