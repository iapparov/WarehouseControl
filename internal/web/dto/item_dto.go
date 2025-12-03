package dto

type ItemCreateRequest struct {
	Name  string  `json:"name" binding:"required"`
	Count int     `json:"count" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type ItemUpdateRequest struct {
	Name  string  `json:"name" binding:"required"`
	Count int     `json:"count" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}
