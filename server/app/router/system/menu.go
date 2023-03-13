package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitMenuRouter(r *gin.RouterGroup) gin.IRouter {
	menuApi := system.NewMenuMenuApi()
	router := r.Group("/menu")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", menuApi.GetMenuTree)
		router.GET("/list", menuApi.GetMenus)
		router.POST("/create", menuApi.CreateMenu)
		router.PATCH("/update/:menuId", menuApi.UpdateMenuById)
		router.DELETE("/delete/batch", menuApi.BatchDeleteMenuByIds)
		router.GET("/access/list/:userId", menuApi.GetUserMenusByUserId)
		router.GET("/access/tree/:userId", menuApi.GetUserMenuTreeByUserId)
	}

	return r
}
