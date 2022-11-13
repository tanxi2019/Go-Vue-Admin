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
	DB *gorm.DB // DB 全局mysql数据库变量

	Redis *redis.Client // Redis全局redis数据库变量

	CasbinEnforcer *casbin.Enforcer // CasbinEnforcer 全局CasbinEnforcer

	Log *zap.SugaredLogger // Log 全局日志变量

	Trans ut.Translator // Trans 全局validate翻译器

	AuthMiddleware *jwt.GinJWTMiddleware // AuthMiddleware jwt auth认证
)
