package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(r *gin.RouterGroup) gin.IRouter {
	api := system.NewApiApi()
	router := r.Group("/api")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())

	{
		router.GET("/list", api.GetApis)
		router.GET("/tree", api.GetApiTree)
		router.POST("/create", api.CreateApi)
		router.PATCH("/update/:apiId", api.UpdateApiById)
		router.DELETE("/delete/batch", api.BatchDeleteApiByIds)
	}

	return r
}
