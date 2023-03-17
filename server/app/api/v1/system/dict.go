package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	model "server/app/model/system"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
)

// DictService
type DictService interface {
	PostDict(c *gin.Context)      // 创建
	GetDictList(c *gin.Context)   // 列表
	PutDict(c *gin.Context)       // 更新
	DeleteDict(c *gin.Context)    // 删除
	DeleteDictAll(c *gin.Context) // 批量删除
}

// DictApiService 服务层数据处理
type DictApiService struct {
	Dict service.DictService
}

// NewDictApi 构造函数
func NewDictApi() DictService {
	dict := service.NewDictService()
	dictService := DictApiService{Dict: dict}
	return dictService
}

// PostDict 创建
func (ds DictApiService) PostDict(c *gin.Context) {
	dict := new(model.Dict)
	// 参数绑定
	if err := c.ShouldBindJSON(&dict); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.Dict.PostDict(dict)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// GetExampleList 列表
func (ds DictApiService) GetDictList(c *gin.Context) {
	var pageList reqo.PageList
	// 参数绑定
	if err := c.ShouldBindQuery(&pageList); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 参数校验

	// 服务层数据操作
	data, total, err := ds.Dict.GetDictList(&pageList)

	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"data":  data,
		"total": total,
	})
	return
}

// PutExample 更新
func (ds DictApiService) PutDict(c *gin.Context) {
	dict := new(model.Dict)
	// 参数绑定
	if err := c.ShouldBind(&dict); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.Dict.PutDict(dict)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteDict  删除
func (ds DictApiService) DeleteDict(c *gin.Context) {
	dictId := new(reqo.DictId)
	if err := c.ShouldBind(&dictId); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.Dict.DeleteDict(dictId)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteDictAll  批量删除
func (ds DictApiService) DeleteDictAll(c *gin.Context) {
	dictIds := new(reqo.DictIds)
	// 参数绑定
	if err := c.ShouldBind(&dictIds); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.Dict.DeleteDictAll(dictIds.DictIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}
