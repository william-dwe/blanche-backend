package dto

type PaginationRequest struct {
	Limit int `json:"limit" form:"limit,default=10" binding:"number,max=100,min=1"`
	Page  int `json:"page" form:"page,default=1" binding:"number"`
}

type PaginationResponse struct {
	TotalData   int64 `json:"total_data"`
	TotalPage   int64 `json:"total_page"`
	CurrentPage int   `json:"current_page"`
}
