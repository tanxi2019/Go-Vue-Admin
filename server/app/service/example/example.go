package example

import (
	"context"
	"github.com/go-redis/redis/v8"
	"server/app/model/example"
	"server/app/model/example/reqo"
	"server/cache"
	"server/global"
)

// ExampleService
type ExampleService interface {
	PostExample(example *example.Example) error                                               // 创建
	GetExample(ID *reqo.ExampleId) (example *example.Example, err error)                      // 单条数据
	GetExampleList(pageInfo *reqo.PageList) (data []*example.Example, total int64, err error) // 列表
	DeleteExample(ID *reqo.ExampleId) error                                                   // 删除
	DeleteExampleAll(exampleIds []uint) error                                                 // 批量删除
	PutExample(example *example.Example) error                                                // 更新
}

type Example struct {
}

// NewExampleService 构造函数
func NewExampleService() ExampleService {
	return Example{}
}

// PostExample 创建
func (ep Example) PostExample(example *example.Example) error {
	err := global.DB.Create(&example).Error
	return err
}

// GetExample 单条数据
func (ep Example) GetExample(ExampleId *reqo.ExampleId) (example *example.Example, err error) {
	// redis 点击数
	global.Redis.Incr(context.Background(), cache.ExampleCountKey(ExampleId.ID))

	// get from cache
	example, err = cache.GetExampleCache(uint64(ExampleId.ID))

	if err == redis.Nil || err != nil {
		// mysql
		err = global.DB.Where("id = ?", ExampleId.ID).First(&example).Error
		_ = cache.SetExampleCache(uint64(ExampleId.ID), example)
		return example, err

	} else {
		// redis
		return example, err
	}

}

// GetExampleList 列表
func (ep Example) GetExampleList(pageInfo *reqo.PageList) (data []*example.Example, total int64, err error) {
	// mysql
	limit := pageInfo.Size
	offset := pageInfo.Size * (pageInfo.Page - 1)
	// 创建db
	db := global.DB.Order("id desc").Model(&example.Example{})
	// 如果有条件搜索 下方会自动创建搜索语句
	if pageInfo.ID != 0 {
		db = db.Where("`id` = ?", int(pageInfo.ID))
	}
	// name
	if pageInfo.Name != "" {
		db = db.Where("`name` LIKE ?", "%"+pageInfo.Name+"%")
	}
	// sex
	if pageInfo.Sex != 0 {
		db = db.Where("`sex` = ?", pageInfo.Sex)
	}
	// age
	if pageInfo.Age > 0 {
		db = db.Where("`age` = ?", pageInfo.Age)
	}
	// mobile
	if pageInfo.Mobile != "" {
		db = db.Where("`mobile` = ?", pageInfo.Mobile)
	}
	// description
	if pageInfo.Description != "" {
		db = db.Where("`description` LIKE ?", "%"+pageInfo.Description+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Find(&data).Error
	_ = cache.SetExampleListCache(data)
	return data, total, err

}

// DeleteExample 删除
func (ep Example) DeleteExample(ID *reqo.ExampleId) error {
	// 软删除,根据主键软删除
	err := global.DB.Delete(&example.Example{}, ID).Error
	return err
}

// DeleteExampleAll 批量删除
func (ep Example) DeleteExampleAll(exampleIds []uint) (err error) {
	// 软删除
	err = global.DB.Where("id IN (?)", exampleIds).Delete(&example.Example{}).Error
	return err
}

// PutExample 更新
func (ep Example) PutExample(example *example.Example) error {
	// 根据id更新
	err := global.DB.Debug().Model(example).Where("id = ?", example.ID).Updates(&example).Error
	return err
}
