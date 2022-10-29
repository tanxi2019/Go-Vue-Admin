package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	model "server/app/model/example"
	"server/global"
	"strconv"
	"time"
)

// ExampleService 用户案例缓存模块接口
type ExampleService interface {
	SetExampleCountCache(id uint64) error                      // set 获取点击数
	GetExampleCountCache(id uint64) uint64                     // get 获取点击数
	SetExampleCache(id uint64, example *model.Example) error   // set 缓存中的单条信息
	GetExampleCache(id uint64) (*model.Example, error)         // get 缓存中的单条信息
	SetExampleRankCache(company []*redis.Z) error              // set 排行榜
	GetExampleRankCache() ([]string, error)                    // get 排行榜
	SetExampleUserVoteCache(id uint, company []*redis.Z) error // set 投票
	GetExampleVote(id uint) ([]string, error)                  // get 投票
	IncrByExampleScore(id uint, t int) error                   // 投票计分,每票5分
	GetExampleScore(id uint) (string, error)                   // 获取投票分数
	ZScoreExample(id uint, member string) (float64, error)     // 获取有序集合存入分数

}

type Example struct{}

// NewExampleService 构造函数
func NewExampleService() ExampleService {
	return Example{}
}

// 点击数
func (e Example) SetExampleCountCache(id uint64) error {
	key := fmt.Sprintf("count:%s", strconv.Itoa(int(id)))
	err := global.Redis.Incr(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}

// 获取点击数
func (e Example) GetExampleCountCache(id uint64) uint64 {
	key := fmt.Sprintf("count:%s", strconv.Itoa(int(id)))
	countStr, _ := global.Redis.Get(context.Background(), key).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// GetExampleCache 获取缓存中的单条信息
func (e Example) GetExampleCache(id uint64) (*model.Example, error) {
	key := fmt.Sprintf("example:%s", strconv.FormatUint(id, 10))
	val, err := global.Redis.Get(context.Background(), key).Result()

	if err == redis.Nil || err != nil {
		return nil, err
	} else {
		example := model.Example{}
		if err := json.Unmarshal([]byte(val), &example); err != nil {
			return nil, err
		}
		return &example, nil
	}
}

// SetExampleCache 缓存单条案例信息key到rides
func (e Example) SetExampleCache(id uint64, example *model.Example) error {
	key := fmt.Sprintf("example:%s", strconv.FormatUint(id, 10))
	content, err := json.Marshal(example)
	if err != nil {
		return err
	}
	errSet := global.Redis.Set(context.Background(), key, content, time.Minute*5).Err()
	if errSet != nil {
		return errSet
	}
	return nil
}

// SetExampleRankCache 排行榜
func (e Example) SetExampleRankCache(company []*redis.Z) error {
	key := fmt.Sprintf("Rank")
	err := global.Redis.ZAdd(context.Background(), key, company...).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetExampleRankCache 获取排行榜
func (e Example) GetExampleRankCache() ([]string, error) {
	key := fmt.Sprintf("Rank")
	rank, err := global.Redis.ZRevRange(context.Background(), key, 0, 9).Result()
	if err != nil {
		return nil, err
	}
	return rank, nil
}

// SetExampleUserVoteCache 投票id
func (e Example) SetExampleUserVoteCache(id uint, company []*redis.Z) error {
	key := fmt.Sprintf("Active:%s", strconv.Itoa(int(id)))
	err := global.Redis.ZAdd(context.Background(), key, company...).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetExampleRankCache 获取排行榜
func (e Example) GetExampleVote(id uint) ([]string, error) {
	key := fmt.Sprintf("Active:%s", strconv.Itoa(int(id)))
	vote, err := global.Redis.ZRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return vote, nil
}

// 投票计分,每天3票，每票5分
func (e Example) IncrByExampleScore(id uint, t int) error {
	key := fmt.Sprintf("Score:%s", strconv.Itoa(int(id)))
	err := global.Redis.IncrBy(context.Background(), key, 5).Err()
	global.Redis.Expire(context.Background(), key, time.Second*time.Duration(t))
	if err != nil {
		return err
	}
	return nil
}

// 获取投票计分
func (e Example) GetExampleScore(id uint) (string, error) {
	key := fmt.Sprintf("Score:%s", strconv.Itoa(int(id)))
	str, err := global.Redis.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}

// 获取有序集合存入分数
func (e Example) ZScoreExample(id uint, member string) (float64, error) {
	key := fmt.Sprintf("Active:%s", strconv.Itoa(int(id)))
	result, err := global.Redis.ZScore(context.Background(), key, member).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
