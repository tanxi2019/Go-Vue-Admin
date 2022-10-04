package reqo

// 字典详情列表结构体
type DictDetailList struct {
	ID      uint   `json:"id" form:"id"`
	KeyWord string `json:"keyword" form:"keyword"`
	DictID  uint   `json:"dict_id" form:"dict_id"`
	Page    int    `json:"page" form:"page"` // 页码
	Size    int    `json:"size" form:"size"` // 每页大小
}

// DictDetaiId 批量删除接口结构体
type DictDetaiId struct {
	ID uint `json:"id" form:"id"`
}

// DictDetaiIds 批量删除接口结构体
type DictDetaiIds struct {
	DictDetaiIds []uint `json:"dictDetailIds" form:"dictDetailIds"`
}
