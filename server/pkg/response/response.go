package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回成功
func Success(c *gin.Context, code int, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code, // 自定义code
		"msg":  msg,  // message
		"data": data, // 数据
	})
	return
}

// 返回失败
func Error(c *gin.Context, httpCode int, code int, msg string, data interface{}) {
	c.JSON(httpCode, map[string]interface{}{
		"code": code, // 自定义code
		"msg":  msg,  // message
		"data": data, // 数据
	})
	return
}
