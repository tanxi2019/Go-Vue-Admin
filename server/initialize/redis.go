package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"server/config"
	"server/global"
)

// InitRedis Redis连接池配置
func InitRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password, // no password set
		DB:       config.Conf.Redis.Database, // use default DB
	})
	global.Redis = redisClient

	result := redisClient.Ping(context.Background())

	global.Log.Infof("初始化redis数据库完成!")

	if result.Val() != "PONG" {
		// 连接有问题
		fmt.Println("Redis 链接失败")
	}

}
