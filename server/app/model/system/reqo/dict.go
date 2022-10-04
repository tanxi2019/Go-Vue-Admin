package reqo

// PageList 分页时使用
type PageList struct {
	ID      uint   `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	KeyWord string `json:"keyword" form:"keyword"`
	Desc    string `json:"desc" form:"desc"`
	Status  *bool  `json:"status" form:"status"`
	Sort    int    `json:"sort" form:"sort"`
	Creator string `json:"creator" form:"creator"`
	Page    int    `json:"page" form:"page"` // 页码
	Size    int    `json:"size" form:"size"` // 每页大小
}

// DictID 批量删除接口结构体
type DictId struct {
	ID uint `json:"id" form:"id"`
}

// DictIds 批量删除接口结构体
type DictIds struct {
	DictIds []uint `json:"DictIds" form:"dictIds"`
}
