package repo

import "server/app/model/system"

// 返回给前端的当前用户信息
type UserInfoResp struct {
	ID           uint           `json:"id"`
	Username     string         `json:"username"`
	Mobile       string         `json:"mobile"`
	Avatar       string         `json:"avatar"`
	Nickname     string         `json:"nickname"`
	Introduction string         `json:"introduction"`
	Roles        []*system.Role `json:"roles"`
}

func ToUserInfoResp(user system.User) UserInfoResp {
	return UserInfoResp{
		ID:           user.ID,
		Username:     user.Username,
		Mobile:       user.Mobile,
		Avatar:       user.Avatar,
		Nickname:     *user.Nickname,
		Introduction: *user.Introduction,
		Roles:        user.Roles,
	}
}

// 返回给前端的用户列表
type UsersResp struct {
	ID           uint   `json:"ID"`
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       uint   `json:"status"`
	Creator      string `json:"creator"`
	RoleIds      []uint `json:"roleIds"`
}

func ToUsersResp(userList []*system.User) []UsersResp {
	var users []UsersResp
	for _, user := range userList {
		userResp := UsersResp{
			ID:           user.ID,
			Username:     user.Username,
			Mobile:       user.Mobile,
			Avatar:       user.Avatar,
			Nickname:     *user.Nickname,
			Introduction: *user.Introduction,
			Status:       user.Status,
			Creator:      user.Creator,
		}
		roleIds := make([]uint, 0)
		for _, role := range user.Roles {
			roleIds = append(roleIds, role.ID)
		}
		userResp.RoleIds = roleIds
		users = append(users, userResp)
	}

	return users
}
