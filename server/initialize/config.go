package initialize

import (
	"fmt"
	"os"
	"server/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 系统配置，对应yml
// viper内置了map structure, yml文件用"-"区分单词, 转为驼峰方便

// InitConfig 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	// 读取配置信息
	err = viper.ReadInConfig()

	// 热更新配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 将读取的配置信息保存至全局变量Conf
		if err := viper.Unmarshal(config.Conf); err != nil {
			panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
		}
	})

	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(config.Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}
}
