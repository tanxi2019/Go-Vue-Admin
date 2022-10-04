package upload

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Aliyun struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"accessKeyId" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"accessKeySecret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucketName" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucketUrl" yaml:"bucket-url"`
	BasePath        string `mapstructure:"base-path" json:"basePath" yaml:"base-path"`
}

type AliyunOss interface {
	UploadFile(folderPath string, file *multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

func NewAliyunOSS(q Aliyun) AliyunOss {
	return &q
}

func (a *Aliyun) UploadFile(folderPath string, file *multipart.FileHeader) (p string, filename string, err error) {
	bucket, err := newBucket(a.Endpoint, a.AccessKeyId, a.AccessKeySecret, a.BucketName)
	if err != nil {
		return "", "", errors.New("function AliyunOss.newBucket() Failed, err:" + err.Error())
	}

	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer f.Close()
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	// yunFileTmpPath := filepath.Join("uploads", time.Now().Format("2006-01-02")) + "/" + file.Filename
	yunFileTmpPath := folderPath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + file.Filename

	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		return "", "", errors.New("function formUploader.Put() Failed, err:" + err.Error())
	}
	p = a.BucketUrl + "/" + yunFileTmpPath
	filename = yunFileTmpPath
	return p, filename, nil
}

func (a *Aliyun) DeleteFile(key string) error {
	bucket, err := newBucket(a.Endpoint, a.AccessKeyId, a.AccessKeySecret, a.BucketName)
	if err != nil {
		return errors.New("function AliyunOss.newBucket() Failed, err:" + err.Error())
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}

	return nil
}

func newBucket(endpoint, accessKeyId, accessKeySecret, bucketName string) (*oss.Bucket, error) {
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}
