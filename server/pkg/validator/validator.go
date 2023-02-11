package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"server/global"
	"server/pkg/code"
	"server/pkg/response"
	"strings"
)

// HandleValidatorError 处理字段校验异常
func HandleValidatorError(c *gin.Context, err error) {
	//如何返回错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	response.Error(c, http.StatusBadRequest, code.ValidateError, code.GetErrMsg(code.ValidateError),
		removeTopStruct(errs.Translate(global.Trans)))
	return
}

//   removeTopStruct 定义一个去掉结构体名称前缀的自定义方法：
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
