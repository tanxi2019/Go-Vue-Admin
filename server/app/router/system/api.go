package system

import (
	"server/app/api/v1/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(r *gin.RouterGroup) gin.IRouter {
	apiApi := system.NewApiApi()
	router := r.Group("/api")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())

	{
		router.GET("/list", apiApi.GetApis)
		router.GET("/tree", apiApi.GetApiTree)
		router.POST("/create", apiApi.CreateApi)
		router.PATCH("/update/:apiId", apiApi.UpdateApiById)
		router.DELETE("/delete/batch", apiApi.BatchDeleteApiByIds)
	}

	return r
}
