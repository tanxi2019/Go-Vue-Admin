package reqo

// PageList 分页时使用
type PageList struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Age         int    `json:"age" form:"age"`
	Sex         int    `json:"sex" form:"sex" `
	Mobile      string `json:"mobile"  form:"mobile"`
	Description string `json:"description" form:"description"  `
	Page        int    `json:"page" form:"page"` // 页码
	Size        int    `json:"size" form:"size"` // 每页大小
}

// ExampleId 批量删除接口结构体
type ExampleId struct {
	ID uint `json:"id" form:"id"`
}

// ExampleId 批量删除接口结构体
type ActiveId struct {
	ID  uint `json:"id" form:"id"`
	VID uint `json:"vid" form:"vid"`
}

// DeleteExample 批量删除接口结构体
type ExampleIds struct {
	ExampleIds []uint `json:"exampleIds" form:"exampleIds"`
}
