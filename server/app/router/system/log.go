package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitOperationLogRouter(r *gin.RouterGroup) gin.IRouter {
	log := system.NewLogApi()
	router := r.Group("/log")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", log.GetOperationLogs)
		router.DELETE("/operation/delete/batch", log.BatchDeleteOperationLogByIds)
	}
	return r
}
