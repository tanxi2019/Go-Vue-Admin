package system

import (
	"github.com/gin-gonic/gin"
	"server/app/api/system"
	"server/global"
)

// InitBaseRouter 注册基础路由
func InitBaseRouter(r *gin.RouterGroup) gin.IRouter {
	base := system.NewCaptchaApi()
	router := r.Group("/base")
	{
		router.POST("/captcha", base.Captcha)
		router.POST("/login", global.AuthMiddleware.LoginHandler)
		router.POST("/logout", global.AuthMiddleware.LogoutHandler)
		router.POST("/refreshToken", global.AuthMiddleware.RefreshHandler)
	}
	return r
}
