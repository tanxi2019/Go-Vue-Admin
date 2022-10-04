package system

import (
	"server/global"

	"github.com/gin-gonic/gin"
)

// InitBaseRouter 注册基础路由
func InitBaseRouter(r *gin.RouterGroup) gin.IRouter {
	router := r.Group("/base")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/login", global.AuthMiddleware.LoginHandler)
		router.POST("/logout", global.AuthMiddleware.LogoutHandler)
		router.POST("/refreshToken", global.AuthMiddleware.RefreshHandler)
	}
	return r
}
