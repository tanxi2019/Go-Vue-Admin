package system

import (
	"net/http"
	"server/app/model/system"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiService interface {
	GetApis(c *gin.Context)             // 获取接口列表
	GetApiTree(c *gin.Context)          // 获取接口树(按接口Category字段分类)
	CreateApi(c *gin.Context)           // 创建接口
	UpdateApiById(c *gin.Context)       // 更新接口
	BatchDeleteApiByIds(c *gin.Context) // 批量删除接口
}

type ApiApiService struct {
	Api service.ApiService
}

func NewApiApi() ApiService {
	return ApiApiService{Api: service.NewApiService()}
}

// GetApis 获取接口列表
func (as ApiApiService) GetApis(c *gin.Context) {
	var req reqo.ApiListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 获取
	data, total, err := as.Api.GetApis(&req)
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

// GetApiTree 获取接口树(按接口Category字段分类)
func (as ApiApiService) GetApiTree(c *gin.Context) {
	tree, err := as.Api.GetApiTree()
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"apiTree": tree,
	})
	return

}

// CreateApi 创建接口
func (as ApiApiService) CreateApi(c *gin.Context) {
	var req reqo.CreateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 获取当前用户
	ur := service.NewUserService()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	api := system.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	// 创建接口
	err = as.Api.CreateApi(&api)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// UpdateApiById 更新接口
func (as ApiApiService) UpdateApiById(c *gin.Context) {
	var req reqo.UpdateApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取路径中的apiId
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	// 获取当前用户
	ur := service.NewUserService()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	api := system.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	err = as.Api.UpdateApiById(uint(apiId), &api)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// BatchDeleteApiByIds 批量删除接口
func (as ApiApiService) BatchDeleteApiByIds(c *gin.Context) {
	var req reqo.DeleteApiRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 删除接口
	err := as.Api.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}
