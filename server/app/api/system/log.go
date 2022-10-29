package system

import (
	"net/http"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"

	"github.com/gin-gonic/gin"
)

type LogApi interface {
	GetOperationLogs(c *gin.Context)             // 获取操作日志列表
	BatchDeleteOperationLogByIds(c *gin.Context) //批量删除操作日志
}

type LogApiService struct {
	Log service.LogService
}

func NewLogApi() LogApi {
	return LogApiService{Log: service.NewLogService()}
}

// GetOperationLogs 获取操作日志列表
func (ls LogApiService) GetOperationLogs(c *gin.Context) {
	var req reqo.OperationLogListRequest
	// 绑定参数
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 获取
	data, total, err := ls.Log.GetOperationLogs(&req)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"data":  data,
		"total": total,
	})
	return
}

// BatchDeleteOperationLogByIds 批量删除操作日志
func (ls LogApiService) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req reqo.DeleteOperationLogRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 删除接口
	err := ls.Log.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}
