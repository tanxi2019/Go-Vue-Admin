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

const (
	// DailyExampleKey 每日排行
	DailyExampleKey = "example:daily"
)

// ExampleCountKey 案例点击数的key
func ExampleCountKey(id uint) string {
	return fmt.Sprintf("count:%s", strconv.Itoa(int(id)))
}

// 缓存单条案例信息key
func getExampletKey(id uint64) string {
	return fmt.Sprintf("example:%s", strconv.FormatUint(id, 10))
}

// 缓存单条案例信息key
func getExampletListKey() string {
	return fmt.Sprintf("example")
}

// 获取点击数
func GetExampleCountCache(id uint64) uint64 {
	countStr, _ := global.Redis.Get(context.Background(), ExampleCountKey(uint(id))).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// GetExampleCache 获取缓存中的单条信息
func GetExampleCache(id uint64) (*model.Example, error) {
	key := getExampletKey(id)
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
func SetExampleCache(id uint64, example *model.Example) error {
	key := getExampletKey(id)
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

// GetExampleListCache 获取缓存中的单条信息
func GetExampleListCache() (example []*model.Example, err error) {
	key := getExampletListKey()
	val, err := global.Redis.Get(context.Background(), key).Result()

	if err == redis.Nil || err != nil {
		return nil, err
	} else {

		if err := json.Unmarshal([]byte(val), &example); err != nil {
			return nil, err
		}
		return example, nil
	}
}

// SetExampleListCache 缓存单条案例信息key到rides
func SetExampleListCache(example []*model.Example) error {
	key := getExampletListKey()
	data, err := json.Marshal(example)
	if err != nil {
		return err
	}
	errSet := global.Redis.Set(context.Background(), key, data, time.Minute*5).Err()
	if errSet != nil {
		return errSet
	}
	return nil
}
