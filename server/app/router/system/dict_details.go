package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/system"
)

// InitDictDetailsRouter 用户案例模块
func InitDictDetailsRouter(r *gin.RouterGroup) gin.IRouter {
	dictDetailsRouter := system.NewDictDetailsApi()
	router := r.Group("/dict/details")
	// 开启jwt认证中间件
	//router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())

	{
		router.POST("/create", dictDetailsRouter.PostDictDetails)        // 创建
		router.GET("/list", dictDetailsRouter.GetDictDetailsList)        // 列表
		router.PUT("/put", dictDetailsRouter.PutDictDetails)             // 更新
		router.DELETE("/delete", dictDetailsRouter.DeleteDictDetails)    // 删除
		router.DELETE("/remove", dictDetailsRouter.DeleteDictDetailsAll) // 全部删除
	}

	return r
}
