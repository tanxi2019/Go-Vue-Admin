package middleware

import (
	"net/http"
	service "server/app/service/system"
	"server/config"
	"server/global"
	"server/pkg/response"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var checkLock sync.Mutex

// CasbinMiddleware Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ur := service.NewUserService()
		user, err := ur.GetCurrentUser(c)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 401, "用户未登录", nil)
			c.Abort()
			return
		}
		if user.Status != 1 {
			response.Error(c, http.StatusBadRequest, 401, "当前用户已被禁用", nil)
			c.Abort()
			return
		}
		// 获得用户的全部角色
		roles := user.Roles
		// 获得用户全部未被禁用的角色的Keyword
		var subs []string
		for _, role := range roles {
			if role.Status == 1 {
				subs = append(subs, role.Keyword)
			}
		}
		// 获得请求路径URL
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// 获取请求方式
		act := c.Request.Method

		isPass := check(subs, obj, act)
		if !isPass {
			response.Error(c, http.StatusBadRequest, 401, "没有权限", nil)

			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()
	isPass := false
	for _, sub := range subs {
		pass, _ := global.CasbinEnforcer.Enforce(sub, obj, act)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}
