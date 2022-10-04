package initialize

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"server/config"
	"server/global"
)

// InitCasbinEnforcer 初始化casbin策略管理器
func InitCasbinEnforcer() {
	enforcer, err := mysqlCasbin()
	if err != nil {
		global.Log.Panicf("初始化Casbin失败：%v", err)
		panic(fmt.Sprintf("初始化Casbin失败：%v", err))
	}

	global.CasbinEnforcer = enforcer
	global.Log.Info("初始化Casbin完成!")
}

func mysqlCasbin() (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		return nil, err
	}
	new, err := casbin.NewEnforcer(config.Conf.Casbin.ModelPath, adapter)
	if err != nil {
		return nil, err
	}

	err = new.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return new, nil
}
