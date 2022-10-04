package example

import (
	"github.com/gin-gonic/gin"
	"server/app/api/example"
	"server/global"
	"server/middleware"
)

// InitExampleRouter 用户案例模块
func InitExampleRouter(r *gin.RouterGroup) gin.IRouter {
	exampleRouter := example.NewExampleApi()
	router := r.Group("/example")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())

	{
		router.POST("/create", exampleRouter.PostExample)        // 创建
		router.GET("/id", exampleRouter.GetExample)              // 单条数据
		router.GET("/list", exampleRouter.GetExampleList)        // 列表
		router.PUT("/put", exampleRouter.PutExample)             // 更新
		router.DELETE("/delete", exampleRouter.DeleteExample)    // 删除
		router.DELETE("/remove", exampleRouter.DeleteExampleAll) // 全部删除
	}

	return r
}
