package logic

import (
	"app/internal/model"
	"app/tools/conv"
	"errors"
)

type AdminAuth struct {
	Id           uint              `json:"id,omitempty"`
	Pid          uint              `json:"pid,omitempty"`
	RolesGroup   *model.RolesGroup `json:"roles_group,omitempty"`
	RolesIds     []uint            `json:"roles_ids,omitempty"`
	IsSuperAdmin bool              `json:"is_super_admin,omitempty"`
	Name         string            `json:"name,omitempty"`
	Avatar       string            `json:"avatar,omitempty"`
}

var AdminList map[uint]*AdminAuth

func init() {
	AdminList = make(map[uint]*AdminAuth)
}

func NewAdminAuth(id, pid uint, RolesGroup *model.RolesGroup, isSuper bool) *AdminAuth {
	group := model.RolesGroup{}
	auth := AdminAuth{}
	auth.Id = id
	auth.Pid = pid
	auth.IsSuperAdmin = isSuper
	if isSuper {
		auth.RolesGroup = nil
		auth.RolesIds = nil
	} else {
		auth.RolesGroup = RolesGroup
		explode, _ := conv.Explode[uint](",", RolesGroup.RolesIds)
		auth.RolesIds = group.GetRolesIdsByIds(explode)
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

// AuthRules 校验权限
func (a *AdminAuth) AuthRules(ids []uint) error {
	if a.IsSuperAdmin {
		return nil
	}
	for _, roleId := range ids {
		if k, _ := conv.InSlice[uint](a.RolesIds, roleId); k == -1 {
			return errors.New("AuthRules Fail")
		}
	}
	return nil
}
