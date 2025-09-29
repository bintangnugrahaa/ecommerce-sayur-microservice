package response

type DefaultResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type DefaultResponseWithPaginations struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination struct {
		Page       int64 `json:"page"`
		TotalCount int64 `json:"total_count"`
		PerPage    int64 `json:"per_page"`
		TotalPage  int64 `json:"total_page"`
	}
}
