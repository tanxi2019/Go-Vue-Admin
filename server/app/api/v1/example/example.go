package example

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"net/http"
	"os"
	model "server/app/model/example"
	"server/app/model/example/reqo"
	service "server/app/service/example"
	"server/cache"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
	"server/serializer"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ExampleService
type ExampleService interface {
	PostExample(c *gin.Context)      // 创建
	GetExample(c *gin.Context)       // 单条数据
	GetExampleList(c *gin.Context)   // 列表
	PutExample(c *gin.Context)       // 更新
	DeleteExample(c *gin.Context)    // 删除
	DeleteExampleAll(c *gin.Context) // 批量删除
	GetExampleRank(c *gin.Context)   // 排行榜
	GetExampleVote(c *gin.Context)   // 投票
	AliPay(c *gin.Context)           // 支付宝
}

// ExampleApiService 服务层数据处理
type ExampleApiService struct {
	Example service.ExampleService
}

// NewExampleApi 创建构造函数工厂模式
func NewExampleApi() ExampleService {
	return ExampleApiService{Example: service.NewExampleService()}
}

// @Tags Example
// @Summary 创建
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Example true "用户ID"
// @Success 200 {object} response.Response{data,msg=string} "获取单一客户信息,返回包括客户详情"
// @Router /api/example/create [post]
func (es ExampleApiService) PostExample(c *gin.Context) {
	example := new(model.Example)
	// 参数绑定
	if err := c.ShouldBindJSON(&example); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
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

// @Tags Example
// @Summary 单条信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query reqo.ExampleId true "用户ID"
// @Success 200 {object} response.Response{data,msg=string} "获取单一客户信息,返回包括客户详情"
// @Router /api/example/id [get]
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

// @Tags Example
// @Summary 分页列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query reqo.PageList true "页码, 每页大小"
// @Success 200 {object} response.Response{data,msg=string} "分页获取权限客户列表,返回包括列表,总数,页码,每页数量"
// @Router /api/example/list [get]
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

// @Tags Example
// @Summary 更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Example true "用户ID"
// @Success 200 {object} response.Response{data,msg=string} "删除客户"
// @Router /api/example/put [put]
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

// @Tags Example
// @Summary 删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query reqo.ExampleId true "用户ID"
// @Success 200 {object} response.Response{msg=string} "删除客户"
// @Router /api/example/delete [delete]
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

// @Tags Example
// @Summary 批量删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query reqo.ExampleIds true "用户ID"
// @Success 200 {object} response.Response{data,msg=string} "删除客户"
// @Router /api/example/remove [delete]
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
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// @Tags Example
// @Summary 排行榜
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data,msg=string}
// @Router /api/example/rank [get]
func (es ExampleApiService) GetExampleRank(c *gin.Context) {
	// 服务层数据操作
	data, err := es.Example.GetExampleRankList()
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 排序：降序
	list := serializer.BuildExampleList(data)
	sort.Slice(list, func(i, j int) bool {
		return list[i].Count >= list[j].Count
	})

	// 成功返回，序列化
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), list)
	return
}

// @Tags Example
// @Summary 投票
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query reqo.ActiveId true "用户ID，活动ID"
// @Success 200 {object} response.Response{data,msg=string}
// @Router /api/example/vote [post]
func (es ExampleApiService) GetExampleVote(c *gin.Context) {
	ActiveId := new(reqo.ActiveId)
	if err := c.ShouldBind(&ActiveId); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	if ActiveId.ID == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "错误", nil)
		return
	}
	ExampleCache := cache.NewExampleService()
	// 投票额度，每票5分
	MAX := 15
	Score, _ := ExampleCache.GetExampleScore(ActiveId.ID)

	// 字符串转 int
	count, _ := strconv.Atoi(Score)
	if count < MAX {
		// 服务层数据操作
		num, err := es.Example.GetExampleVote(ActiveId)
		if err != nil {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
			return
		}
		// 获取所有投票id
		vote, err := ExampleCache.GetExampleVote(ActiveId.VID)
		if err != nil {
			return
		}
		// 获取所有投票分数
		var arr []float64
		for i := 0; i < len(vote); i++ {
			// 获取分数
			result, err := ExampleCache.ZScoreExample(ActiveId.VID, vote[i])
			if err != nil {
				response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
			}

			arr = append(arr, result)
		}
		// 定义一个管道
		ch := make(chan float64)
		// 开启协程
		go es.Sum(arr[:len(arr)/2], ch)
		go es.Sum(arr[len(arr)/2:], ch)
		// 从管道中取出数据
		a, b := <-ch, <-ch
		sum := a + b

		// 成功返回，序列化
		response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
			"num":   num / 5, // 序列化
			"total": sum,
		})
		return
	} else {
		response.Error(c, http.StatusBadRequest, code.ServerErr, "今日投票已达上限", nil)
		return
	}

}

// 计数数字之和
func (es ExampleApiService) Sum(arr []float64, ch chan float64) {
	var result float64 = 0
	for _, v := range arr {
		result += v
	}
	// 向管道中传输数据
	ch <- result
}

// @Tags Example
// @Summary  支付宝支付接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data,msg=string}
// @Router /api/example/alipay [get]
func (es ExampleApiService) AliPay(c *gin.Context) {
	// 公钥
	aliPublicKey, _ := os.ReadFile("./config/alipay/publicKey.txt")
	// 私钥
	aliPrivateKey, _ := os.ReadFile("./config/alipay/privateKey.txt")
	// appid
	appId := "2021000117699350"

	client, err := alipay.NewClient(appId, string(aliPrivateKey), false)
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 自定义配置http请求接收返回结果body大小，默认 10MB
	client.SetBodySize(10) // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认关闭
	//client.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	// 注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).               // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2).              // 设置签名类型，不设置默认 RSA2
							SetReturnUrl("https://www.baidu.com"). // 设置返回URL
							SetNotifyUrl("https://www.baidu.com")  // 设置异步通知URL
	//SetAppAuthToken("")               // 设置第三方应用授权
	// 公钥验签
	client.AutoVerifySign([]byte(aliPublicKey))
	// 订单号
	trade := strings.Replace(time.Now().Format("2006 01 02 15 04 05"), " ", "", -1)
	// 初始化结构体
	bm := make(gopay.BodyMap)
	bm.Set("subject", "电脑网站测试支付")                    // 商品名称
	bm.Set("out_trade_no", trade)                    // 订单号
	bm.Set("total_amount", "100.00")                 // 金额
	bm.Set("product_code", "FAST_INSTANT_TRADE_PAY") // FAST_INSTANT_TRADE_PAY电脑支付 QUICK_WAP_WAY 手机支付
	// 返回支付链接地址
	payUrl, err := client.TradePagePay(c, bm)
	// 失败返回
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返沪
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"url": payUrl,
	})
	return

}
