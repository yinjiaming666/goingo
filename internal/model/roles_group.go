package model

import (
	"app/tools/conv"
)

// RolesGroup 角色表
type RolesGroup struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint   `json:"id" gorm:"primaryKey; type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Name            string `json:"name" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Pid             uint   `json:"pid" gorm:"type:INT(11) UNSIGNED NOT NULL; default:0"`
	RolesIds        string `json:"roles_ids" gorm:"type:VARCHAR(1200) NOT NULL; default:''; comment: 权限id(逗号分隔)"`
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
