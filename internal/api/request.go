package api

type PaginationQueryParam struct {
	Size int `form:"size" binding:"min=1,max=20"`
	Page int `form:"page" binding:"min=0"`
}
