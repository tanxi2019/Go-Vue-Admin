package captcha

import (
	"time"

	"github.com/mojocn/base64Captcha"
)

//CaptchaReq
type CaptchaReq struct {
	ImgHeight int `json:"imgHeight"`
	ImgWidth  int `json:"imgWidth"`
	KeyLong   int `json:"keyLong"`
}

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
}

// 1.开辟一个验证码使用的存储空间 3000个 时间3分钟内有效
var store = base64Captcha.NewMemoryStore(3000, 3*time.Minute)

// 生成验证码
func GenCaptcha(ca CaptchaReq) (result CaptchaResponse, err error) {
	// 2.创建验证码驱动 五种：dight 数字验证码；audio 语音验证码；string 字符验证码；math 数学验证码(加减乘除)；chinese中文验证码-有bug
	// 图片高度80 宽度240 数字位数6 最大绝对偏斜因子 背景圆圈数量
	driver := base64Captcha.NewDriverDigit(ca.ImgHeight, ca.ImgWidth, ca.KeyLong, 0.7, 80)
	// 3.生成验证码并保存至store
	cp := base64Captcha.NewCaptcha(driver, store)
	// 4.生成base64图像及id
	id, b64s, err := cp.Generate()
	if err != nil {
		return result, err
	}
	result = CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: ca.KeyLong,
	}
	return result, nil
}

// 验证验证码 原验证码的id，待验证的输入字符串answer
func VerifyCaptcha(id string, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return store.Verify(id, answer, true)
}
