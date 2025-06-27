package api

type PaginationQueryParam struct {
	Size uint `form:"size" binding:"min=1,max=20"`
	Page uint `form:"page" binding:"min=0"`
}
