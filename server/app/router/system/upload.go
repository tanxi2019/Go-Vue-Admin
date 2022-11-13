package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/system"
)

// InitUploadRouter 用户案例模块
func InitUploadRouter(r *gin.RouterGroup) gin.IRouter {
	uploadRouter := system.NewUploadApi()
	router := r.Group("/upload")
	////开启jwt认证中间件
	//router.Use(global.AuthMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/file", uploadRouter.UploadFile)      // 创建
		router.POST("/qiniu", uploadRouter.UploadQiniuYun) // 创建
	}

	return r
}
