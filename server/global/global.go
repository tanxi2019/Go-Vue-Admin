package global

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/casbin/casbin/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// DB 全局mysql数据库变量
	DB *gorm.DB

	// Redis全局redis数据库变量
	Redis *redis.Client

	// CasbinEnforcer 全局CasbinEnforcer
	CasbinEnforcer *casbin.Enforcer

	// Log 全局日志变量
	Log *zap.SugaredLogger

	// Trans 全局validate翻译器
	Trans ut.Translator

	// AuthMiddleware jwt auth认证
	AuthMiddleware *jwt.GinJWTMiddleware
)
