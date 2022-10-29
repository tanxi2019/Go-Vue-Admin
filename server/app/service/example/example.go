package example

import (
	"github.com/go-redis/redis/v8"
	"server/app/model/example"
	"server/app/model/example/reqo"
	"server/cache"
	"server/global"
	"strconv"
	"time"
)

// ExampleService
type ExampleService interface {
	PostExample(example *example.Example) error                                               // 创建
	GetExample(ID *reqo.ExampleId) (example *example.Example, err error)                      // 单条数据
	GetExampleList(pageInfo *reqo.PageList) (data []*example.Example, total int64, err error) // 列表
	DeleteExample(ID *reqo.ExampleId) error                                                   // 删除
	DeleteExampleAll(exampleIds []uint) error                                                 // 批量删除
	PutExample(example *example.Example) error                                                // 更新
	GetExampleRankList() (data []*example.Example, err error)                                 // 排行榜
	GetExampleVote(ActiveId *reqo.ActiveId) (int, error)                                      // 投票
}

type Example struct{}

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
	ExampleCache := cache.NewExampleService()
	err = ExampleCache.SetExampleCountCache(uint64(ExampleId.ID))
	if err != nil {
		return nil, err
	}
	// 排行榜
	company := []*redis.Z{
		{Score: float64(ExampleCache.GetExampleCountCache(uint64(ExampleId.ID))), Member: ExampleId.ID},
	}
	err = ExampleCache.SetExampleRankCache(company)
	if err != nil {
		return nil, err
	}

	// get from cache
	example, err = ExampleCache.GetExampleCache(uint64(ExampleId.ID))

	if err == redis.Nil || err != nil {
		// mysql
		err = global.DB.Where("id = ?", ExampleId.ID).First(&example).Error
		_ = ExampleCache.SetExampleCache(uint64(ExampleId.ID), example)
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
	err := global.DB.Model(example).Where("id = ?", example.ID).Updates(&example).Error
	return err
}

// GetExampleRankList 排行榜列表
func (ep Example) GetExampleRankList() (data []*example.Example, err error) {
	// redis
	ExampleCache := cache.NewExampleService()
	rank, err := ExampleCache.GetExampleRankCache()
	if err != nil {
		return nil, err
	}
	// mysql
	err = global.DB.Model(&example.Example{}).Where("id in ?", rank).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err

}

// GetExample 投票
func (ep Example) GetExampleVote(ActiveId *reqo.ActiveId) (int, error) {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	// 24 h
	day := 60 * 60 * 24
	// 当前时分秒
	timer := hour*60*60 + minute*60 + second
	// 剩余过期时间
	t := day - timer
	ExampleCache := cache.NewExampleService()
	// 投票计数,每票5分
	err := ExampleCache.IncrByExampleScore(ActiveId.ID, t)
	if err != nil {
		return 0, err
	}

	// 判断是否存在
	Score, err := ExampleCache.GetExampleScore(ActiveId.ID)
	if err != nil {
		return 0, err
	}
	// 字符串转 int
	score, _ := strconv.Atoi(Score)
	// 用户投票记录
	company := []*redis.Z{
		{Score: float64(score), Member: ActiveId.ID},
	}
	err = ExampleCache.SetExampleUserVoteCache(ActiveId.VID, company)
	if err != nil {
		return 0, err
	}
	return score, nil

}
