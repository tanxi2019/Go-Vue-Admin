package system

import (
	"server/app/api/v1/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// InitUserRouter 注册用户路由
func InitUserRouter(r *gin.RouterGroup) gin.IRouter {
	userApi := system.NewUserApi()
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/info", userApi.GetUserInfo)
		router.GET("/list", userApi.GetUsers)
		router.PUT("/changePwd", userApi.ChangePwd)
		router.POST("/create", userApi.CreateUser)
		router.PATCH("/update/:userId", userApi.UpdateUserById)
		router.DELETE("/delete/batch", userApi.BatchDeleteUserByIds)
	}
	return r
}
