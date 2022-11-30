package system

import (
	"errors"
	"fmt"
	"server/app/model/system"
	"server/app/model/system/reqo"
	"server/global"
	"strings"
)

type RoleService interface {
	GetRoles(req *reqo.RoleListRequest) ([]system.Role, int64, error)    // 获取角色列表
	GetRolesByIds(roleIds []uint) ([]*system.Role, error)                // 根据角色ID获取角色
	CreateRole(role *system.Role) error                                  // 创建角色
	UpdateRoleById(roleId uint, role *system.Role) error                 // 更新角色
	GetRoleMenusById(roleId uint) ([]*system.Menu, error)                // 获取角色的权限菜单
	UpdateRoleMenus(role *system.Role) error                             // 更新角色的权限菜单
	GetRoleApisByRoleKeyword(roleKeyword string) ([]*system.Api, error)  // 根据角色关键字获取角色的权限接口
	UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error // 更新角色的权限接口（先全部删除再新增）
	BatchDeleteRoleByIds(roleIds []uint) error                           // 删除角色
}

type Role struct{}

func NewRoleService() RoleService {
	return Role{}
}

// GetRoles 获取角色列表
func (rd Role) GetRoles(req *reqo.RoleListRequest) ([]system.Role, int64, error) {
	var list []system.Role
	db := global.DB.Model(&system.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// GetRolesByIds 根据角色ID获取角色
func (rd Role) GetRolesByIds(roleIds []uint) ([]*system.Role, error) {
	var list []*system.Role
	err := global.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// CreateRole 创建角色
func (rd Role) CreateRole(role *system.Role) error {
	err := global.DB.Create(role).Error
	return err
}

// UpdateRoleById 更新角色
func (rd Role) UpdateRoleById(roleId uint, role *system.Role) error {
	err := global.DB.Model(&system.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// GetRoleMenusById 获取角色的权限菜单
func (rd Role) GetRoleMenusById(roleId uint) ([]*system.Menu, error) {
	var role system.Role
	err := global.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// UpdateRoleMenus 更新角色的权限菜单
func (rd Role) UpdateRoleMenus(role *system.Role) error {
	err := global.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// GetRoleApisByRoleKeyword 根据角色关键字获取角色的权限接口
func (rd Role) GetRoleApisByRoleKeyword(roleKeyword string) ([]*system.Api, error) {
	policies := global.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	// 获取所有接口
	var apis []*system.Api
	err := global.DB.Find(&apis).Error
	if err != nil {
		return apis, errors.New("获取角色的权限接口失败")
	}

	accessApis := make([]*system.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// UpdateRoleApis 更新角色的权限接口（先全部删除再新增）
func (rd Role) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := global.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies := global.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := global.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	isAdded, _ := global.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("更新角色的权限接口失败")
	}
	err = global.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}

// BatchDeleteRoleByIds 删除角色
func (rd Role) BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*system.Role
	err := global.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = global.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := global.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := global.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}

	}
	return err
}
