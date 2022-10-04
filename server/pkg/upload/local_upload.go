package upload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

type Local struct{}

type LocalOss interface {
	UploadFile(folderPath string, file *multipart.FileHeader) (string, string, error)
	DeleteFile(folderPath string, key string) error
}

func NewLocalOss() LocalOss {
	return &Local{}
}

// 存储文件
// folderPath文件夹根目录 file文件对象
// p文件路径 filename文件名 err错误
func (*Local) UploadFile(folderPath string, file *multipart.FileHeader) (p string, filename string, err error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = MD5V([]byte(name))
	// 拼接新文件名
	filename = name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(folderPath, os.ModePerm)
	if mkdirErr != nil {
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p = folderPath + "/" + filename

	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()

	out, createErr := os.Create(p)
	if createErr != nil {
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()

	_, copyErr := io.Copy(out, f)
	if copyErr != nil {
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}

// 删除文件
// folderPath文件夹根目录 key文件名
func (*Local) DeleteFile(folderPath string, key string) error {
	p := folderPath + "/" + key
	if strings.Contains(p, folderPath) {
		if err := os.Remove(p); err != nil {
			return errors.New("delete local file error, err:" + err.Error())
		}
	}
	return nil
}

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
