package example

import (
	"github.com/gin-gonic/gin"
	"server/app/api/example"
	"server/global"
	"server/middleware"
)

// InitExampleRouter 用户案例模块
func InitExampleRouter(r *gin.RouterGroup) gin.IRouter {
	// api层
	exampleApi := example.NewExampleApi()
	// 路由组
	router := r.Group("/example")
	//开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/create", exampleApi.PostExample)        // 创建
		router.GET("/id", exampleApi.GetExample)              // 单条数据
		router.GET("/list", exampleApi.GetExampleList)        // 列表
		router.PUT("/put", exampleApi.PutExample)             // 更新
		router.DELETE("/delete", exampleApi.DeleteExample)    // 删除
		router.DELETE("/remove", exampleApi.DeleteExampleAll) // 全部删除
		router.GET("/rank", exampleApi.GetExampleRank)        // 排行榜
		router.POST("/vote", exampleApi.GetExampleVote)       // 投票
	}

	return r
}
