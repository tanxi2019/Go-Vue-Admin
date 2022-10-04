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

// DictDetailsService
type DictDetailsService interface {
	PostDictDetails(c *gin.Context)      // 创建
	GetDictDetailsList(c *gin.Context)   // 列表
	PutDictDetails(c *gin.Context)       // 更新
	DeleteDictDetails(c *gin.Context)    // 删除
	DeleteDictDetailsAll(c *gin.Context) // 批量删除
}

// DictDetailsApiService 服务层数据处理
type DictDetailsApiService struct {
	DictDetails service.DictDetailsService
}

// NewDictDetailsApi 构造函数
func NewDictDetailsApi() DictDetailsService {
	dictDetails := service.NewDictDetailsService()
	dictDetailsApiService := DictDetailsApiService{DictDetails: dictDetails}
	return dictDetailsApiService
}

// PostDictDetails 创建
func (ds DictDetailsApiService) PostDictDetails(c *gin.Context) {
	example := new(model.DictDetail)
	// 参数绑定
	if err := c.ShouldBindJSON(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.DictDetails.PostDictDetails(example)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// GetDictDetailsList 列表
func (ds DictDetailsApiService) GetDictDetailsList(c *gin.Context) {

	var pageList reqo.DictDetailList
	// 参数绑定
	if err := c.ShouldBindQuery(&pageList); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 参数校验

	// 服务层数据操作
	data, total, err := ds.DictDetails.GetDictDetailsList(&pageList)

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

// PutDictDetails 更新
func (ds DictDetailsApiService) PutDictDetails(c *gin.Context) {
	example := new(model.DictDetail)
	// 参数绑定
	if err := c.ShouldBind(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.DictDetails.PutDictDetails(example)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteDictDetails 删除
func (ds DictDetailsApiService) DeleteDictDetails(c *gin.Context) {
	DictId := new(reqo.DictId)
	if err := c.ShouldBind(&DictId); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.DictDetails.DeleteDictDetails(DictId)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteDictDetailsAll 批量删除
func (ds DictDetailsApiService) DeleteDictDetailsAll(c *gin.Context) {
	dictDetaiIds := new(reqo.DictDetaiIds)
	// 参数绑定
	if err := c.ShouldBind(&dictDetaiIds); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := ds.DictDetails.DeleteDictDetailsAll(dictDetaiIds.DictDetaiIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}
