package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	service "server/app/service/system"
	"server/pkg/code"
	"server/pkg/response"
	"strconv"
	"time"
)

type UploadService interface {
	UploadFile(c *gin.Context)     // 本地文件
	UploadQiniuYun(c *gin.Context) // 七牛云文件
}

// UploadApiService 服务层数据处理
type UploadApiService struct {
	Upload service.UploadService
}

// NewUploadApi 创建构造函数简单工厂模式
func NewUploadApi() UploadService {
	return UploadApiService{Upload: service.NewQiniuYunService()}
}

// 本地文件上传
func (u UploadApiService) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	err = c.SaveUploadedFile(file, "./static/img/"+strconv.FormatInt(time.Now().Unix(), 10)+file.Filename)
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), file)
	return

}

// 七牛云文件上传
func (u UploadApiService) UploadQiniuYun(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	url, err := u.Upload.UploadQiniuYun(file)
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), url)
	return

}
