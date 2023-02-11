package router

import (
	"fmt"
	ginSwagger "github.com/swaggo/gin-swagger"
	"server/app/api/socket"
	"server/app/router/example"
	"server/app/router/system"
	"server/config"
	_ "server/docs"
	"server/global"
	"server/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
)

// InitRouter 初始化
func InitRouter() *gin.Engine {
	//设置模式
	gin.SetMode(config.Conf.System.Mode)
	r := gin.Default()

	// 启用限流中间件
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	quantum := config.Conf.RateLimit.Quantum

	// 每100毫秒产生quantum个令牌，桶容量是capacity，1s=1000ms，QPS=quantum*10
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity, quantum))

	// 启用全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 启用操作日志中间件
	r.Use(middleware.OperationLogMiddleware())

	// 初始化JWT认证中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		global.Log.Panicf("初始化JWT中间件失败：%v", err)
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}
	global.AuthMiddleware = authMiddleware

	GroupRouter(r)
	// r.Static(relativePath string, root string)
	global.Log.Info("初始化路由完成！")

	return r
}

func GroupRouter(r *gin.Engine) {
	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// websocket
	r.GET("/ws", socket.Handler)
	// 路由分组
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)
	// 注册路由
	system.InitBaseRouter(apiGroup)         // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	system.InitUserRouter(apiGroup)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	system.InitRoleRouter(apiGroup)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	system.InitMenuRouter(apiGroup)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	system.InitApiRouter(apiGroup)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	system.InitOperationLogRouter(apiGroup) // 注册日志路由, jwt认证中间件,casbin鉴权中间件
	system.InitDictRouter(apiGroup)         // 注册字典路由, jwt认证中间件,casbin鉴权中间件
	system.InitDictDetailsRouter(apiGroup)  // 注册字典详情路由, jwt认证中间件,casbin鉴权中间件
	system.InitUploadRouter(apiGroup)       // 文件上传, jwt认证中间件,casbin鉴权中间件
	example.InitExampleRouter(apiGroup)     // 注册基础路由, jwt认证中间件,要casbin中间件

}
