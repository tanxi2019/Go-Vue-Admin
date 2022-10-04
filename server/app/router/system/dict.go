package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/system"
)

// InitExampleRouter 用户案例模块
func InitDictRouter(r *gin.RouterGroup) gin.IRouter {
	dictRouter := system.NewDictApi()
	router := r.Group("/dict")
	// 开启jwt认证中间件
	//router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())

	{
		router.POST("/create", dictRouter.PostDict)        // 创建
		router.GET("/list", dictRouter.GetDictList)        // 列表
		router.PUT("/put", dictRouter.PutDict)             // 更新
		router.DELETE("/delete", dictRouter.DeleteDict)    // 删除
		router.DELETE("/remove", dictRouter.DeleteDictAll) // 全部删除
	}

	return r
}
