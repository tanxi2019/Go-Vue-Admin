package initialize

import (
	"server/config"

	"github.com/gogf/gf/frame/g"
)

func InitializeI18n() {
	g.I18n().SetLanguage(config.Conf.System.I18nLanguage)
}

/**
i18n国际化和validate trans中的翻译：
都应该根据请求头中的Accept-Language来决定设置成什么语言，
配置文件的设置方式只能实现在不同国家的服务器部署的时候去选择官方语言
**/
