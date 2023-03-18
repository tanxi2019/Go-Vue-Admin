package system

import (
	"server/app/api/v1/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(r *gin.RouterGroup) gin.IRouter {
	roleApi := system.NewRoleApi()
	router := r.Group("/role")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", roleApi.GetRoles)
		router.POST("/create", roleApi.CreateRole)
		router.PATCH("/update/:roleId", roleApi.UpdateRoleById)
		router.GET("/menus/get/:roleId", roleApi.GetRoleMenusById)
		router.PATCH("/menus/update/:roleId", roleApi.UpdateRoleMenusById)
		router.GET("/apis/get/:roleId", roleApi.GetRoleApisById)
		router.PATCH("/apis/update/:roleId", roleApi.UpdateRoleApisById)
		router.DELETE("/delete/batch", roleApi.BatchDeleteRoleByIds)
	}
	return r
}
