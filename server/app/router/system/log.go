package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitOperationLogRouter(r *gin.RouterGroup) gin.IRouter {
	logApi := system.NewLogApi()
	router := r.Group("/log")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", logApi.GetOperationLogs)
		router.DELETE("/operation/delete/batch", logApi.BatchDeleteOperationLogByIds)
	}
	return r
}
