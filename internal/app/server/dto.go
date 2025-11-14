package server

type errorResponseDTO struct {
	Error errorDTO `json:"error"`
}

type errorDTO struct {
	Message string `json:"message"`
}

type saveResultDTO struct {
	TotalCount      int     `json:"total_count"`
	DuplicatesCount int     `json:"duplicates_count"`
	TotalItems      int     `json:"total_items"`
	TotalCategories int     `json:"total_categories"`
	TotalPrice      float32 `json:"total_price"`
}
