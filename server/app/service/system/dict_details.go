package system

import (
	"server/app/model/system"
	"server/app/model/system/reqo"
	"server/global"
)

// DictDetailsService
type DictDetailsService interface {
	PostDictDetails(dictDetail *system.DictDetail) error                                                  // 创建
	GetDictDetailsList(pageInfo *reqo.DictDetailList) (data []*system.DictDetail, total int64, err error) // 列表
	DeleteDictDetails(ID *reqo.DictId) error                                                              // 删除
	DeleteDictDetailsAll(dictDetailsIds []uint) error                                                     // 批量删除
	PutDictDetails(dictDetail *system.DictDetail) error                                                   // 更新
}

type DictDetails struct {
}

// NewDictDetailsService 构造函数
func NewDictDetailsService() DictDetailsService {
	return DictDetails{}
}

// PostDictDetails 创建
func (ds DictDetails) PostDictDetails(dictDetail *system.DictDetail) error {
	err := global.DB.Create(&dictDetail).Error
	return err
}

// GetDictDetailsList 列表
func (ds DictDetails) GetDictDetailsList(pageInfo *reqo.DictDetailList) (data []*system.DictDetail, total int64, err error) {
	// gorm 获获列表数据
	limit := pageInfo.Size
	offset := pageInfo.Size * (pageInfo.Page - 1)

	// sql语句
	sql := "SELECT * FROM dict_detail WHERE dict_id IN ( " +
		"SELECT id FROM dict WHERE key_word = ? and status IS True  ) and deleted_at IS NULL order by sort limit ?,?;"

	if err = global.DB.Raw(sql, pageInfo.KeyWord, offset, limit).Scan(&data).Error; err != nil {
		return data, 0, err
	}

	db := global.DB.Debug().Order("sort asc").Model(&system.DictDetail{})

	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return data, total, err
}

// DeleteDictDetails 删除
func (ds DictDetails) DeleteDictDetails(ID *reqo.DictId) error {
	// 软删除,根据主键软删除
	err := global.DB.Delete(&system.DictDetail{}, ID).Error
	return err
}

// DeleteExampleAll 批量删除
func (ds DictDetails) DeleteDictDetailsAll(dictDetailsIds []uint) (err error) {
	// 软删除
	err = global.DB.Where("id IN (?)", dictDetailsIds).Delete(&system.DictDetail{}).Error
	return err
}

// PutExample 更新
func (ds DictDetails) PutDictDetails(dictDetail *system.DictDetail) error {
	// 根据id更新
	err := global.DB.Model(dictDetail).Where("id = ?", dictDetail.ID).Updates(&dictDetail).Error
	return err
}
