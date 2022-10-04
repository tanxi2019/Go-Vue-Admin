package reqo

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page int `json:"page" form:"page"` // 页码
	Size int `json:"size" form:"size"` // 每页大小
}
