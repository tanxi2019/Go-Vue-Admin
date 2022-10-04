package system

import (
	"net/http"
	"server/app/model/system"
	"server/app/model/system/repo"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/pkg/bcrypt"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type UserApi interface {
	GetUserInfo(c *gin.Context)          // 获取当前登录用户信息
	GetUsers(c *gin.Context)             // 获取用户列表
	ChangePwd(c *gin.Context)            // 更新用户登录密码
	CreateUser(c *gin.Context)           // 创建用户
	UpdateUserById(c *gin.Context)       // 更新用户
	BatchDeleteUserByIds(c *gin.Context) // 批量删除用户
}

// UserService 服务层数据处理
type UserService struct {
	User service.UserService
}

// NewUserApi 构造函数
func NewUserApi() UserApi {
	user := service.NewUserService()
	userService := UserService{User: user}
	return userService
}

// GetUserInfo 获取当前登录用户信息
func (us UserService) GetUserInfo(c *gin.Context) {
	user, err := us.User.GetCurrentUser(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	userInforesponsep := repo.ToUserInfoResp(user)
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"userInfo": userInforesponsep,
	})
	return

}

// GetUsers 获取用户列表
func (us UserService) GetUsers(c *gin.Context) {
	var req reqo.UserListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 获取
	users, total, err := us.User.GetUsers(&req)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"users": repo.ToUsersResp(users),
		"total": total,
	})
	return
}

// ChangePwd 更新用户登录密码
func (us UserService) ChangePwd(c *gin.Context) {
	var req reqo.ChangePwdRequest

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取当前用户
	user, err := us.User.GetCurrentUser(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = bcrypt.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "原密码有误", nil)
		return
	}
	// 更新密码
	err = us.User.ChangePwd(user.Username, bcrypt.GenPasswd(req.NewPassword))
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "更新密码失败", nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// CreateUser 创建用户
func (us UserService) CreateUser(c *gin.Context) {
	var req reqo.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	if len(req.Password) < 6 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "密码长度至少为6位", nil)
		return

	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := us.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "根据角色ID获取角色信息失败: "+err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "用户不能创建比自己等级高的或者相同等级的用户", nil)
		return
	}

	// 密码为空就默认 123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := system.User{
		Username:     req.Username,
		Password:     bcrypt.GenPasswd(req.Password),
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}

	err = us.User.CreateUser(&user)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, code.GetErrMsg(code.ServerErr), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}

// UpdateUserById 更新用户
func (us UserService) UpdateUserById(c *gin.Context) {
	var req reqo.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	//获取path中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "用户ID不正确", nil)
		return
	}

	// 根据path中的userId获取用户信息
	oldUser, err := us.User.GetUserById(uint(userId))
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取需要更新的用户信息失败: "+err.Error(), nil)
		return

	}

	// 获取当前用户
	ctxUser, err := us.User.GetCurrentUser(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "根据角色ID获取角色信息失败: "+err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息: "+err.Error(), nil)
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := system.User{
		Model:        oldUser.Model,
		Username:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if userId == int(ctxUser.ID) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "不能禁用自己", nil)
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "不能更改自己的角色", nil)
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "请到个人中心更新自身密码", nil)
			return
		}

		// 密码赋值
		user.Password = ctxUser.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := us.User.GetUserMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "根据用户ID获取用户角色排序最小值失败", nil)
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "用户不能更新比自己角色等级高的或者相同等级的用户", nil)
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "用户不能把别的用户角色等级更新得比自己高或相等", nil)
			return
		}

		// 密码赋值
		if req.Password != "" {
			user.Password = bcrypt.GenPasswd(req.Password)
		}

	}

	// 更新用户
	err = us.User.UpdateUser(&user)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "更新用户失败", nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}

// BatchDeleteUserByIds 批量删除用户
func (us UserService) BatchDeleteUserByIds(c *gin.Context) {
	var req reqo.DeleteUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 前端传来的用户ID
	reqUserIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := us.User.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "根据用户ID获取用户角色排序最小值失败", nil)
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := us.User.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqUserIds, ctxUser.ID) {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "用户不能删除自己", nil)
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "用户不能删除比自己角色等级高的用户", nil)
			return
		}
	}

	err = us.User.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "删除用户失败: "+err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}
