package system

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"server/config"
	"strconv"
	"time"
)

// UploadService
type UploadService interface {
	UploadQiniuYun(file *multipart.FileHeader) (string, error) // 七牛云文件
}

type Upload struct{}

// NewExampleService 构造函数
func NewQiniuYunService() UploadService {
	return Upload{}
}

// UploadQiniuYun 七牛云文件
func (u Upload) UploadQiniuYun(file *multipart.FileHeader) (string, error) {
	src, _ := file.Open()
	defer src.Close()
	putPolicy := storage.PutPolicy{
		Scope: config.Conf.Qiniu.Bucket,
	}
	// 秘钥验证
	mac := qbox.NewMac(config.Conf.Qiniu.AccessKey, config.Conf.Qiniu.SecretKey)
	fmt.Println(config.Conf.Qiniu.AccessKey)
	fmt.Println(config.Conf.Qiniu.SecretKey)
	// 上传凭证
	Token := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong, // 华东
		UseCdnDomains: false,
		UseHTTPS:      false, // not https
	}

	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 上传路径，如果当前目录中已存在相同文件
	key := "/image/" + strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	err := formUploader.Put(context.Background(), &ret, Token, key, src, file.Size, &putExtra)
	if err != nil {
		return "", err
	}

	url := config.Conf.Qiniu.Origin + ret.Key
	return url, nil
}
