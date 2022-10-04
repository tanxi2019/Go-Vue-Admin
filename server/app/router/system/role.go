package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(r *gin.RouterGroup) gin.IRouter {
	role := system.NewRoleApi()
	router := r.Group("/role")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", role.GetRoles)
		router.POST("/create", role.CreateRole)
		router.PATCH("/update/:roleId", role.UpdateRoleById)
		router.GET("/menus/get/:roleId", role.GetRoleMenusById)
		router.PATCH("/menus/update/:roleId", role.UpdateRoleMenusById)
		router.GET("/apis/get/:roleId", role.GetRoleApisById)
		router.PATCH("/apis/update/:roleId", role.UpdateRoleApisById)
		router.DELETE("/delete/batch", role.BatchDeleteRoleByIds)
	}
	return r
}
