package example

import (
	"gorm.io/gorm"
)

// Example 结构体
type Example struct {
	gorm.Model
	Name        string `json:"name" form:"name" binding:"required" gorm:"type:varchar(20);not null;unique" `
	Age         int    `json:"age" form:"age" gorm:"type:int(3);not null"`
	Sex         int    `json:"sex" form:"sex" gorm:"type:int(2);not null"`
	Mobile      string `json:"mobile" form:"mobile" gorm:"type:varchar(11);not null" `
	Description string `json:"description" form:"description" gorm:"type:varchar(100);not null" `
}
