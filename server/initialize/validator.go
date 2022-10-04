package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"server/config"
	"server/global"
	"strings"
)

// InitTrans validator信息翻译
func InitValidate() (err error) {
	var Trans ut.Translator
	//修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		Trans, ok = uni.GetTranslator(config.Conf.System.I18nLanguage)
		global.Trans = Trans
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", config.Conf.System.I18nLanguage)
		}
		switch config.Conf.System.I18nLanguage {
		case "en":
			_ = en_translations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			_ = zh_translations.RegisterDefaultTranslations(v, Trans)
		default:
			_ = en_translations.RegisterDefaultTranslations(v, Trans)
		}

		return
	}
	return
}
