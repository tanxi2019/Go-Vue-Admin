package example

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	model "server/app/model/example"
	"server/app/model/example/reqo"
	service "server/app/service/example"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
	"server/serializer"
)

//ExampleService
type ExampleService interface {
	PostExample(c *gin.Context)      // 创建
	GetExample(c *gin.Context)       // 单条数据
	GetExampleList(c *gin.Context)   // 列表
	PutExample(c *gin.Context)       // 更新
	DeleteExample(c *gin.Context)    // 删除
	DeleteExampleAll(c *gin.Context) // 批量删除
}

// ExampleApiService 服务层数据处理
type ExampleApiService struct {
	Example service.ExampleService
}

// NewExampleApi 构造函数
func NewExampleApi() ExampleService {
	example := service.NewExampleService()
	exampleService := ExampleApiService{Example: example}
	return exampleService
}

// PostExample 创建
func (es ExampleApiService) PostExample(c *gin.Context) {
	example := new(model.Example)
	// 参数绑定
	if err := c.ShouldBindJSON(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	fmt.Println(example)
	// 服务层数据操作
	err := es.Example.PostExample(example)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// GetExample 单条数据
func (es ExampleApiService) GetExample(c *gin.Context) {
	ExampleId := new(reqo.ExampleId)
	if err := c.ShouldBindQuery(&ExampleId); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 服务层数据操作
	data, err := es.Example.GetExample(ExampleId)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回，序列化
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), serializer.BuildExample(data))
	return
}

// GetExampleList 列表
func (es ExampleApiService) GetExampleList(c *gin.Context) {
	var pageList reqo.PageList
	// 参数绑定
	if err := c.ShouldBindQuery(&pageList); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 服务层数据操作
	data, total, err := es.Example.GetExampleList(&pageList)

	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"data":  serializer.BuildExampleList(data), // 序列化
		"total": total,
	})
	return
}

// PutExample 更新
func (es ExampleApiService) PutExample(c *gin.Context) {
	example := new(model.Example)
	// 参数绑定
	if err := c.ShouldBind(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := es.Example.PutExample(example)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteExample 删除
func (es ExampleApiService) DeleteExample(c *gin.Context) {
	DeleteExampleID := new(reqo.ExampleId)
	if err := c.ShouldBind(&DeleteExampleID); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := es.Example.DeleteExample(DeleteExampleID)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// DeleteExampleAll 批量删除
func (es ExampleApiService) DeleteExampleAll(c *gin.Context) {
	example := new(reqo.ExampleIds)
	// 参数绑定
	if err := c.ShouldBind(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 服务层数据操作
	err := es.Example.DeleteExampleAll(example.ExampleIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}
