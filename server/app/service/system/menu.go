package system

import (
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"server/app/model/system"
	"server/global"
)

type MenuService interface {
	GetMenus() ([]*system.Menu, error)                   // 获取菜单列表
	GetMenuTree() ([]*system.Menu, error)                // 获取菜单树
	CreateMenu(menu *system.Menu) error                  // 创建菜单
	UpdateMenuById(menuId uint, menu *system.Menu) error // 更新菜单
	BatchDeleteMenuByIds(menuIds []uint) error           // 批量删除菜单

	GetUserMenusByUserId(userId uint) ([]*system.Menu, error)    // 根据用户ID获取用户的权限(可访问)菜单列表
	GetUserMenuTreeByUserId(userId uint) ([]*system.Menu, error) // 根据用户ID获取用户的权限(可访问)菜单树
}

type Menu struct{}

func NewMenuService() MenuService {
	return Menu{}
}

// GetMenus 获取菜单列表
func (m Menu) GetMenus() ([]*system.Menu, error) {
	var menus []*system.Menu
	err := global.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// GetMenuTree 获取菜单树
func (m Menu) GetMenuTree() ([]*system.Menu, error) {
	var menus []*system.Menu
	err := global.DB.Order("sort").Find(&menus).Error
	// parentId为0的是根菜单
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentId uint, menus []*system.Menu) []*system.Menu {
	tree := make([]*system.Menu, 0)

	for _, md := range menus {
		if *md.ParentId == parentId {
			children := GenMenuTree(md.ID, menus)
			md.Children = children
			tree = append(tree, md)
		}
	}
	return tree
}

// CreateMenu 创建菜单
func (m Menu) CreateMenu(menu *system.Menu) error {
	err := global.DB.Create(menu).Error
	return err
}

// UpdateMenuById 更新菜单
func (m Menu) UpdateMenuById(menuId uint, menu *system.Menu) error {
	err := global.DB.Model(menu).Where("id = ?", menuId).Updates(menu).Error
	return err
}

// BatchDeleteMenuByIds 批量删除菜单
func (m Menu) BatchDeleteMenuByIds(menuIds []uint) error {
	var menus []*system.Menu
	err := global.DB.Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = global.DB.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// GetUserMenusByUserId 根据用户ID获取用户的权限(可访问)菜单列表
func (m Menu) GetUserMenusByUserId(userId uint) ([]*system.Menu, error) {
	// 获取用户
	var user system.User
	err := global.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	allRoleMenus := make([]*system.Menu, 0)
	for _, role := range roles {
		var userRole system.Role
		err := global.DB.Where("id = ?", role.ID).Preload("Menus",
			// 修复排序
			func(db *gorm.DB) *gorm.DB {
				return db.Order("sort")
			}).First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*system.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*system.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// GetUserMenuTreeByUserId 根据用户ID获取用户的权限(可访问)菜单树
func (m Menu) GetUserMenuTreeByUserId(userId uint) ([]*system.Menu, error) {
	menus, err := m.GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}
