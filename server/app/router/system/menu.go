package system

import (
	"server/app/api/system"
	"server/global"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func InitMenuRouter(r *gin.RouterGroup) gin.IRouter {
	menu := system.NewMenuMenuApi()
	router := r.Group("/menu")
	// 开启jwt认证中间件
	router.Use(global.AuthMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", menu.GetMenuTree)
		router.GET("/list", menu.GetMenus)
		router.POST("/create", menu.CreateMenu)
		router.PATCH("/update/:menuId", menu.UpdateMenuById)
		router.DELETE("/delete/batch", menu.BatchDeleteMenuByIds)
		router.GET("/access/list/:userId", menu.GetUserMenusByUserId)
		router.GET("/access/tree/:userId", menu.GetUserMenuTreeByUserId)
	}

	return r
}
