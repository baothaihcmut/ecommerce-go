package pagination

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalPage   int `json:"total_page"`
	TotalItem   int `json:"total_item"`
}

type PaginationResult[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagnination"`
}
