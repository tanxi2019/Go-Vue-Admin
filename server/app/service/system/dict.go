package system

import (
	"server/app/model/system"
	"server/app/model/system/reqo"
	"server/global"
)

// ExampleService
type DictService interface {
	PostDict(dict *system.Dict) error                                                  // 创建
	GetDictList(pageInfo *reqo.PageList) (data []*system.Dict, total int64, err error) // 列表
	DeleteDict(id *reqo.DictId) error                                                  // 删除
	DeleteDictAll(dictIds []uint) error                                                // 批量删除
	PutDict(dict *system.Dict) error                                                   // 更新
}

type Dict struct {
}

// NewDictService 构造函数
func NewDictService() DictService {
	return Dict{}
}

// PostDict 创建
func (dt Dict) PostDict(dict *system.Dict) error {
	err := global.DB.Create(&dict).Error
	return err
}

// GetDictList 列表
func (dt Dict) GetDictList(pageInfo *reqo.PageList) (data []*system.Dict, total int64, err error) {
	// gorm 获获列表数据
	limit := pageInfo.Size
	offset := pageInfo.Size * (pageInfo.Page - 1)
	// 创建db
	db := global.DB.Debug().Order("id desc").Model(&system.Dict{})
	// 如果有条件搜索 下方会自动创建搜索语句
	if pageInfo.ID != 0 {
		db = db.Where("`id` = ?", pageInfo.ID)
	}
	// name
	if pageInfo.Name != "" {
		db = db.Where("`name` LIKE ?", "%"+pageInfo.Name+"%")
	}
	// keyword
	if pageInfo.KeyWord != "" {
		db = db.Where("`key_word` = ?", pageInfo.KeyWord)
	}
	// desc
	if pageInfo.Desc != "" {
		db = db.Where("`desc` LIKE ?", "%"+pageInfo.Desc+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(limit).Offset(offset).Find(&data).Error

	return data, total, err
}

// DeleteDict 删除
func (dt Dict) DeleteDict(id *reqo.DictId) error {
	// 软删除,根据主键软删除
	err := global.DB.Delete(&system.Dict{}, id).Error
	return err
}

// DeleteDictAll 批量删除
func (dt Dict) DeleteDictAll(dictIds []uint) (err error) {
	// 软删除
	err = global.DB.Where("id IN (?)", dictIds).Delete(&system.Dict{}).Error
	return err
}

// PutDict 更新
func (dt Dict) PutDict(dict *system.Dict) error {
	// 根据id更新
	err := global.DB.Model(dict).Where("id = ?", dict.ID).Updates(&dict).Error
	return err
}
