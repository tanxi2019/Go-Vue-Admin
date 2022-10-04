package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// InitUserRouter 注册用户路由
func InitUserRouter(r *gin.RouterGroup) gin.IRouter {
	user := system.NewUserApi()
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.POST("/info", user.GetUserInfo)
		router.GET("/list", user.GetUsers)
		router.PUT("/changePwd", user.ChangePwd)
		router.POST("/create", user.CreateUser)
		router.PATCH("/update/:userId", user.UpdateUserById)
		router.DELETE("/delete/batch", user.BatchDeleteUserByIds)
	}
	return r
}
