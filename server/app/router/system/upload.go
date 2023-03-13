package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/system"
)

// InitUploadRouter 用户案例模块
func InitUploadRouter(r *gin.RouterGroup) gin.IRouter {
	uploadApi := system.NewUploadApi()
	router := r.Group("/upload")
	////开启jwt认证中间件
	//router.Use(global.AuthMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/file", uploadApi.UploadFile)      // 创建
		router.POST("/qiniu", uploadApi.UploadQiniuYun) // 创建
	}

	return r
}
