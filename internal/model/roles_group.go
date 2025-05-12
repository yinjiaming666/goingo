package model

import (
	"app/tools/conv"
)

// RolesGroup 角色表
type RolesGroup struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey; type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	AdminId         uint           `json:"admin_id" gorm:"type:INT(11) UNSIGNED NOT NULL; default:0; comment: 创建人id"`
	RolesIds        string         `json:"roles_ids" gorm:"type:VARCHAR(1200) NOT NULL; default:''; comment: 权限id(逗号分隔)"`
	Pid             uint           `json:"pid" gorm:"type:INT(8) UNSIGNED NOT NULL; default:0; comment: 上级角色组"`
}

type RolesGroupIds struct {
	Id uint `json:"id" gorm:"primaryKey; type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
}

func (r *RolesGroup) GetRolesGroup() {
	Db().Where(r).First(&r)
}

func (r *RolesGroup) GetRolesIdsByIds(ids []uint) []uint {
	list := make([]*RolesGroup, 0)
	Db().Model(r).Where("id in ?", ids).Find(&list)

	uints := make([]uint, 0)
	for _, v := range list {
		explode, _ := conv.Explode[uint](",", v.RolesIds)
		uints = append(uints, explode...)
	}
	return uints
}

func (r *RolesGroup) GetRolesGroupList(adminId uint) []*RolesGroup {
	list := make([]*RolesGroup, 0)
	where := make(map[string]interface{})
	if adminId > 0 {
		where["admin_id"] = adminId
	}
	Db().Where(where).Order("id desc").Find(&list)
	return list
}

func (r *RolesGroup) SetRolesGroup() *RolesGroup {
	if r.Id <= 0 {
		Db().Create(&r)
	} else {
		Db().Model(&r).Updates(&r)
	}
	return r
}

func (r *RolesGroup) DelRolesGroup() *RolesGroup {
	Db().Delete(&r, r.Id)
	return r
}
