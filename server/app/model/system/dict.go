package system

import "gorm.io/gorm"

// 字典分类结构体
type Dict struct {
	gorm.Model
	Name    string `gorm:"type:varchar(50);comment:'字典类别';not null;unique" json:"name" form:"name" binding:"required" validate:"unique" `
	KeyWord string `gorm:"type:varchar(50);comment:'关键词'" json:"keyword"`
	Desc    string `gorm:"type:varchar(100);comment:'说明'" json:"desc"`
	Status  *bool  `gorm:"column:status;default:true; comment:状态 0 1" json:"status"`
	Sort    int    `gorm:"type:varchar(50);comment:'排序'" json:"sort"`
	Creator string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}

// DictCategory结构体表名称
func (Dict) TableName() string {
	return "dict"
}
