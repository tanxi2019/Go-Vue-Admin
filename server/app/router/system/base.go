package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/global"
	"server/pkg/captcha"
	"server/pkg/code"
	"server/pkg/response"
)

// InitBaseRouter 注册基础路由
func InitBaseRouter(r *gin.RouterGroup) gin.IRouter {
	router := r.Group("/base")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/captcha", func(c *gin.Context) {
			CaptchaReq := &captcha.CaptchaReq{
				ImgHeight: 80,
				ImgWidth:  270,
				KeyLong:   6,
			}

			data, err := captcha.GenCaptcha(CaptchaReq)
			if err != nil {
				response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
			}
			response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), data)
			return

		})
		router.POST("/login", global.AuthMiddleware.LoginHandler)
		router.POST("/logout", global.AuthMiddleware.LogoutHandler)
		router.POST("/refreshToken", global.AuthMiddleware.RefreshHandler)
	}
	return r
}
