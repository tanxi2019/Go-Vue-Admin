package serializer

import (
	"server/app/model/example"
	"server/cache"
	"time"
)

// Example 列化器
type Example struct {
	ID          uint      `json:"ID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	DeletedAt   time.Time `json:"DeletedAt"`
	Name        string    `json:"name" `
	Age         int       `json:"age" `
	Sex         int       `json:"sex"  `
	Mobile      string    `json:"mobile"  `
	Count       uint      `json:"count"  `
	Description string    `json:"description"`
}

// BuildExample 序列化
func BuildExample(item *example.Example) Example {
	return Example{
		ID:          item.ID,
		Name:        item.Name,
		Sex:         item.Sex,
		Age:         item.Age,
		Mobile:      item.Mobile,
		Count:       uint(cache.GetExampleCountCache(uint64(item.ID))), // 点击数
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

// BuildExampleList 序列化列表
func BuildExampleList(list []*example.Example) (data []Example) {
	for _, item := range list {
		buildExample := BuildExample(item)
		data = append(data, buildExample)
	}
	return data
}
