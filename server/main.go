package main

import (
	"server/initialize"
)

func main() {
	initialize.InitConfig()          // 加载配置文件到全局配置结构体
	initialize.InitializeI18n()      // 初始化I18N国际化
	initialize.InitLogger()          // 初始化日志
	initialize.InitRedis()           //  初始化数据库(redis)
	initialize.InitMysql()           // 初始化数据库(mysql)
	initialize.InitCasbinEnforcer()  // 初始化casbin策略管理器
	initialize.InitValidate()        // 初始化Validator数据校验
	initialize.InitializeRunServer() // 初始化服务器
}
