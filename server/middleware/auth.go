package middleware

import (
	"net/http"
	"server/app/model/system"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/config"
	"server/global"
	"server/pkg/json"

	"server/pkg/response"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// InitAuth 初始化jwt中间件
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Conf.Jwt.Realm,                                 // jwt标识
		Key:             []byte(config.Conf.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(config.Conf.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(config.Conf.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                           // 有效载荷处理
		IdentityHandler: identityHandler,                                       // 解析Claims
		Authenticator:   login,                                                 // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                          // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                          // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                         // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                        // 登出后的响应
		RefreshResponse: refreshResponse,                                       // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",    // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                              // header名称
		TimeFunc:        time.Now,
	})
	return authMiddleware, err
}

// 有效载荷处理
func payloadFunc(data interface{}) (claims jwt.MapClaims) {
	if v, ok := data.(map[string]interface{}); ok {
		var user system.User
		// 将用户json转为结构体
		err := json.JsonI2Struct(v["user"], &user)
		if err != nil {
			global.Log.Errorw("JsonI2Struct error: ", err)
		} else {
			claims = jwt.MapClaims{
				jwt.IdentityKey: user.ID,
				"user":          v["user"],
			}
		}
	}
	return claims
}

// 解析Claims
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// 校验token的正确性, 处理登录逻辑
func login(c *gin.Context) (interface{}, error) {
	var req reqo.RegisterAndLoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}
	// 密码校验
	userDao := service.NewUserService()
	user, err := userDao.Login(&req)
	if err != nil {
		return nil, err
	}

	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	if userstr, err := json.Struct2Json(user); err != nil {
		global.Log.Errorw("struct2json error: ", err)
		return "", err
	} else {
		return map[string]interface{}{
			"user": userstr,
		}, nil
	}
}

// 用户登录校验成功处理
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userStr := v["user"].(string)
		var user system.User
		// 将用户json转为结构体
		if err := json.Json2Struct(userStr, &user); err != nil {
			global.Log.Errorw("json2struct error: ", err)
			return false
		} else {
			// 将用户保存到context, api调用时取数据方便
			c.Set("user", user)
			return true
		}
	}
	return false
}

// 用户登录校验失败处理
func unauthorized(c *gin.Context, code int, message string) {
	global.Log.Debugf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message)
	// 错误返回
	response.Error(c, http.StatusBadRequest, 10010, message, nil)
	return
}

// 登录成功后的响应
func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	// 成功返回
	response.Success(c, 200, "登录成功", map[string]interface{}{
		"token":   token,
		"expires": expires.Format("2006-01-02 15:04:05"),
	})
	return
}

// 登出后的响应
func logoutResponse(c *gin.Context, code int) {
	// 成功返回
	response.Success(c, 200, "退出成功", nil)
	return
}

// 刷新token后的响应
func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	// 成功返回
	response.Success(c, 200, "刷新token成功", map[string]interface{}{
		"token":   token,
		"expires": expires,
	})
	return
}

/**
jwt json web token 身份提供者和服务提供者间传递被认证的用户身份信息

传统session的验证：
在服务器存储一份用户登录的信息，这份登录信息会在响应时传递给浏览器，告诉其保存为cookie,以便下次请求时发送给我们的应用，这样我们的应用就能识别请求来自哪个用户了；随着不同客户端用户的增加，独立的服务器已无法承载更多的用户

基于token的验证
客户端存储token，并在每次请求时附送上这个token值；服务端验证token值，并返回数据
不需要在服务端去保留用户的认证信息或者会话信息

jwt构成：
头部（header),载荷（payload)，第三部分是签证
**/
