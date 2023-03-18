package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/v1/system"
)

// InitDictDetailsRouter 用户案例模块
func InitDictDetailsRouter(r *gin.RouterGroup) gin.IRouter {
	dictDetailsApi := system.NewDictDetailsApi()
	router := r.Group("/dict/details")
	// 开启jwt认证中间件
	//router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())

	{
		router.POST("/create", dictDetailsApi.PostDictDetails)        // 创建
		router.GET("/list", dictDetailsApi.GetDictDetailsList)        // 列表
		router.PUT("/put", dictDetailsApi.PutDictDetails)             // 更新
		router.DELETE("/delete", dictDetailsApi.DeleteDictDetails)    // 删除
		router.DELETE("/remove", dictDetailsApi.DeleteDictDetailsAll) // 全部删除
	}

	return r
}
