package logic

import (
	"app/internal/model"
	"app/tools/conv"
)

type AdminAuth struct {
	Id            uint
	Pid           uint
	RolesGroupIds []uint
	RolesIds      []uint
	IsSuperAdmin  bool
	Name          string
	Avatar        string
}

var AdminList map[uint]*AdminAuth

func init() {
	AdminList = make(map[uint]*AdminAuth)
}

func NewAdminAuth(id, pid uint, rolesGroupIds string, isSuper bool) *AdminAuth {
	group := model.RolesGroup{}
	auth := AdminAuth{}
	auth.Id = id
	auth.Pid = pid
	auth.IsSuperAdmin = isSuper
	if rolesGroupIds == "*" {
		auth.RolesGroupIds = nil
		auth.RolesIds = nil
	} else {
		auth.RolesGroupIds, _ = conv.Explode[uint](",", rolesGroupIds)
		auth.RolesIds = group.GetRolesIdsByIds(auth.RolesGroupIds)
	}
	return &auth
}

// Cache 缓存权限信息
func (a *AdminAuth) Cache() {
	AdminList[a.Id] = a
}

// ClearCache 清除缓存
func (a *AdminAuth) ClearCache() {
	delete(AdminList, a.Id)
}

func GetAdminAuth(id uint) *AdminAuth {
	return AdminList[id]
}

func (a *AdminAuth) GetAllSonAdmin() {

}

// GetAllRules 获取管理员所有的权限
func (a *AdminAuth) GetAllRules(t int) []*model.RolesFormat {
	roles := model.Roles{}
	var ids = make([]uint, 0)
	if !a.IsSuperAdmin {
		ids = a.RolesIds
	}
	getRoles := roles.GetRoles(ids, t)
	return roles.FormatTree(getRoles)
}
