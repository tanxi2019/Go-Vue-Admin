package code

const (
	SUCCESS       = 200
	ERROR         = 500
	ServerErr     = 1000
	ValidateError = 1001
	Deadline      = 1002
	CreateError   = 1003
	FindError     = 1004
	WithoutServer = 1005
	AuthError     = 1006
	DeleteError   = 1007
	EmptyFile     = 1008
	RateLimit     = 1009
	Unauthorized  = 10010
	WithoutLogin  = 10011
	DisableAuth   = 10012
)

var codeMsg = map[int]string{
	SUCCESS:       "成功",
	ERROR:         "失败",
	ServerErr:     "服务器错误",
	ValidateError: "参数校验错误",
	Deadline:      "服务调用超时",
	CreateError:   "服务器写入失败",
	FindError:     "服务器查询失败",
	WithoutServer: "服务未启用",
	AuthError:     "权限错误",
	DeleteError:   "服务器删除失败",
	EmptyFile:     "文件为空",
	RateLimit:     "访问限流",
	Unauthorized:  "JWT认证失败",
	WithoutLogin:  "用户未登录",
	DisableAuth:   "当前用户已被禁用",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
