package system

import (
	"fmt"
	"net/http"
	"server/app/model/system"
	"server/app/model/system/reqo"
	service "server/app/service/system"
	"server/global"
	"server/pkg/code"
	"server/pkg/response"
	"server/pkg/validator"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type RoleApi interface {
	GetRoles(c *gin.Context)             // 获取角色列表
	CreateRole(c *gin.Context)           // 创建角色
	UpdateRoleById(c *gin.Context)       // 更新角色
	GetRoleMenusById(c *gin.Context)     // 获取角色的权限菜单
	UpdateRoleMenusById(c *gin.Context)  // 更新角色的权限菜单
	GetRoleApisById(c *gin.Context)      // 获取角色的权限接口
	UpdateRoleApisById(c *gin.Context)   // 更新角色的权限接口
	BatchDeleteRoleByIds(c *gin.Context) // 批量删除角色
}

type RoleApiService struct {
	Role service.RoleService
}

func NewRoleApi() RoleApi {
	return RoleApiService{Role: service.NewRoleService()}
}

// GetRoles 获取角色列表
func (rs RoleApiService) GetRoles(c *gin.Context) {
	var req reqo.RoleListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取角色列表
	roles, total, err := rs.Role.GetRoles(&req)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取角色列表失败: "+err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"roles": roles,
		"total": total,
	})
	return
}

// CreateRole 创建角色
func (rs RoleApiService) CreateRole(c *gin.Context) {
	var req reqo.CreateRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取当前用户最高角色等级
	uc := service.NewUserService()
	sort, ctxUser, err := uc.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取当前用户最高角色等级失败: "+err.Error(), nil)
		return
	}

	// 用户不能创建比自己等级高或相同等级的角色
	if sort >= req.Sort {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "不能创建比自己等级高或相同等级的角色", nil)
		return
	}

	role := system.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 创建角色
	err = rs.Role.CreateRole(&role)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "创建角色失败: "+err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}

// UpdateRoleById 更新角色
func (rs RoleApiService) UpdateRoleById(c *gin.Context) {
	var req reqo.CreateRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "角色ID不正确", nil)
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := service.NewUserService()
	_, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	// 不能更新比自己角色等级高或相等的角色
	// 根据path中的角色ID获取该角色信息
	roles, err := rs.Role.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}

	role := system.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// 更新角色
	err = rs.Role.UpdateRoleById(uint(roleId), &role)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色失败: "+err.Error(), nil)
		return
	}

	// 如果更新成功，且更新了角色的keyword, 则更新casbin中policy
	if req.Keyword != roles[0].Keyword {
		// 获取policy
		rolePolicies := global.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色成功", nil)
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// 替换keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}

		// 这里需要先新增再删除（先删除再增加会出错）
		isAdded, _ := global.CasbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色成功，但角色关键字关联的权限接口更新失败", nil)
			return
		}
		isRemoved, _ := global.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色成功，但角色关键字关联的权限接口更新失败", nil)
			return
		}
		err := global.CasbinEnforcer.LoadPolicy()
		if err != nil {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色成功，但角色关键字关联角色的权限接口策略加载失败", nil)
			return
		}

	}

	// 更新角色成功处理用户信息缓存有两种做法:（这里使用第二种方法，因为一个角色下用户数量可能很多，第二种方法可以分散数据库压力）
	// 1.可以帮助用户更新拥有该角色的用户信息缓存,使用下面方法
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2.直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearUserInfoCache()

	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return
}

// GetRoleMenusById 获取角色的权限菜单
func (rs RoleApiService) GetRoleMenusById(c *gin.Context) {
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "角色ID不正确", nil)
		return
	}
	menus, err := rs.Role.GetRoleMenusById(uint(roleId))
	fmt.Println(menus[0])
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取角色的权限菜单失败: "+err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"menus": menus,
	})
	return
}

// UpdateRoleMenusById 更新角色的权限菜单
func (rs RoleApiService) UpdateRoleMenusById(c *gin.Context) {
	var req reqo.UpdateRoleMenusRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "角色ID不正确", nil)
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rs.Role.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := service.NewUserService()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "不能更新比自己角色等级高或相等角色的权限菜单", nil)
			return
		}
	}

	// 获取当前用户所拥有的权限菜单
	mr := service.NewMenuService()
	ctxUserMenus, err := mr.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取当前用户的可访问菜单列表失败: "+err.Error(), nil)
		return
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 前端传来最新的MenuIds集合
	menuIds := req.MenuIds

	// 用户需要修改的菜单集合
	reqMenus := make([]*system.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				// 错误返回
				response.Error(c, http.StatusBadRequest, code.ServerErr, fmt.Sprintf("无权设置ID为%d的菜单", id), nil)
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := mr.GetMenus()
		if err != nil {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "获取菜单列表失败: "+err.Error(), nil)
			return
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = rs.Role.UpdateRoleMenus(roles[0])
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "更新角色的权限菜单失败: "+err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}

// GetRoleApisById 获取角色的权限接口
func (rs RoleApiService) GetRoleApisById(c *gin.Context) {
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "角色ID不正确: ", nil)
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rs.Role.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}
	// 根据角色keyword获取casbin中policy
	keyword := roles[0].Keyword
	apis, err := rs.Role.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), map[string]interface{}{
		"apis": apis,
	})
	return
}

// UpdateRoleApisById 更新角色的权限接口
func (rs RoleApiService) UpdateRoleApisById(c *gin.Context) {
	var req reqo.UpdateRoleApisRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "角色ID不正确", nil)
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rs.Role.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := service.NewUserService()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限接口
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "不能更新比自己角色等级高或相等角色的权限接口", nil)
			return
		}
	}

	// 获取当前用户所拥有的权限接口
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := global.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// 得到path中的角色ID对应角色能够设置的权限接口集合
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// 前端传来最新的ApiID集合
	apiIds := req.ApiIds
	// 根据apiID获取接口详情
	ar := service.NewApiService()
	apis, err := ar.GetApisById(apiIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "根据接口ID获取接口信息失败", nil)
		return
	}
	// 生成前端想要设置的角色policies
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// (非管理员)不能把角色的权限接口设置的比当前用户所拥有的权限接口多
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				// 错误返回
				response.Error(c, http.StatusBadRequest, code.ServerErr, fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]), nil)
				return
			}
		}
	}

	// 更新角色的权限接口
	err = rs.Role.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}

// BatchDeleteRoleByIds 批量删除角色
func (rs RoleApiService) BatchDeleteRoleByIds(c *gin.Context) {
	var req reqo.DeleteRoleRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		// 参数校验
		validator.HandleValidatorError(c, err)
		return
	}

	// 获取当前用户最高等级角色
	ur := service.NewUserService()
	minSort, _, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, err.Error(), nil)
		return
	}

	// 前端传来需要删除的角色ID
	roleIds := req.RoleIds
	// 获取角色信息
	roles, err := rs.Role.GetRolesByIds(roleIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "获取角色信息失败: "+err.Error(), nil)
		return
	}
	if len(roles) == 0 {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "未获取到角色信息", nil)
		return
	}

	// 不能删除比自己角色等级高或相等的角色
	for _, role := range roles {
		if minSort >= role.Sort {
			// 错误返回
			response.Error(c, http.StatusBadRequest, code.ServerErr, "不能删除比自己角色等级高或相等的角色", nil)
			return
		}
	}

	// 删除角色
	err = rs.Role.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		// 错误返回
		response.Error(c, http.StatusBadRequest, code.ServerErr, "删除角色失败", nil)
		return

	}

	// 删除角色成功直接清理缓存，让活跃的用户自己重新缓存最新用户信息
	ur.ClearUserInfoCache()
	// 成功返回
	response.Success(c, code.SUCCESS, code.GetErrMsg(code.SUCCESS), nil)
	return

}
