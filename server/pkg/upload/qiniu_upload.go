package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

type Qiniu struct {
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`                                // 存储区域
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                          // 空间名称
	ImgPath       string `mapstructure:"img-path" json:"imgPath" yaml:"img-path"`                     // CDN加速域名
	UseHTTPS      bool   `mapstructure:"use-https" json:"useHttps" yaml:"use-https"`                  // 是否使用https
	AccessKey     string `mapstructure:"access-key" json:"accessKey" yaml:"access-key"`               // 秘钥AK
	SecretKey     string `mapstructure:"secret-key" json:"secretKey" yaml:"secret-key"`               // 秘钥SK
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"useCdnDomains" yaml:"use-cdn-domains"` // 上传是否使用CDN上传加速
}

type QiniuOss interface {
	UploadFile(folderPath string, file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewQiniuOss(q Qiniu) QiniuOss {
	return &q
}

// 上传文件
func (q *Qiniu) UploadFile(folderPath string, file *multipart.FileHeader) (p string, filekey string, err error) {
	cfg := qiniuConfig(q.UseHTTPS, q.UseCdnDomains, q.Zone)

	formUploader := storage.NewFormUploader(cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()

	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename) // 文件名格式 自己可以改 建议保证唯一性
	putPolicy := storage.PutPolicy{Scope: q.Bucket}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
	}
	p = folderPath + "/" + ret.Key
	filekey = ret.Key
	return p, filekey, nil
}

// 删除文件
func (q *Qiniu) DeleteFile(key string) error {
	cfg := qiniuConfig(q.UseHTTPS, q.UseCdnDomains, q.Zone)

	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	bucketManager := storage.NewBucketManager(mac, cfg)
	if err := bucketManager.Delete(q.Bucket, key); err != nil {
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

// 据配置文件进行返回七牛云的配置
func qiniuConfig(useHTTPS bool, useCdnDomains bool, zone string) *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      useHTTPS,
		UseCdnDomains: useCdnDomains,
	}
	switch zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &cfg
}
