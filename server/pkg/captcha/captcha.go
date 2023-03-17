package captcha

import (
	"github.com/mojocn/base64Captcha"
	"server/cache"
)

// CaptchaReq
type CaptchaReq struct {
	ImgHeight int `json:"imgHeight"`
	ImgWidth  int `json:"imgWidth"`
	KeyLong   int `json:"keyLong"`
}

type CaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
	Captcha   string `json:"captcha"`
}

// 生成验证码
func GenCaptcha(ca *CaptchaReq) (result CaptchaResponse, err error) {
	// 2.创建验证码驱动 五种：dight 数字验证码；audio 语音验证码；string 字符验证码；math 数学验证码(加减乘除)；chinese中文验证码-有bug
	// 图片高度80 宽度240 数字位数6 最大绝对偏斜因子 背景圆圈数量
	driver := base64Captcha.NewDriverDigit(ca.ImgHeight, ca.ImgWidth, ca.KeyLong, 0.7, 80)
	//base64Captcha.NewDriverMath()
	// 3.生成验证码并保存至store
	cp := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	// 4.生成base64图像及id
	id, b64s, err := cp.Generate()
	if err != nil {
		return result, err
	}
	captcha := base64Captcha.DefaultMemStore.Get(id, true)

	// 5. 数字验证码存redis
	CaptchaCache := cache.NewCaptchaService()
	_ = CaptchaCache.SetCaptcha(id, captcha)
	// 返回值
	result = CaptchaResponse{
		CaptchaId: id,
		PicPath:   b64s,
		Captcha:   captcha,
	}

	return result, nil
}

// VerifyCaptcha 验证验证码 原验证码的id，待验证的输入字符串answer
func VerifyCaptcha(id string, answer string) bool {
	return base64Captcha.DefaultMemStore.Verify(id, answer, true)
}
