package model

// AdminRolesGroup 管理员-角色 关联表
type AdminRolesGroup struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint `json:"id" gorm:"primaryKey; type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	AdminId         uint `json:"admin_id" gorm:"type:INT(11) UNSIGNED NOT NULL; default:0; comment: 管理员id"`
	RolesGroupId    uint `json:"roles_group_id" gorm:"type:INT(11) UNSIGNED NOT NULL; default:''; comment: 角色组id"`
}

func (a *AdminRolesGroup) GetGroupIdsByAdminId(adminId uint) []uint {
	list := make([]*RolesGroupIds, 0)
	Db().Model(a).Where("id = ?", adminId).Find(&list)

	uints := make([]uint, 0)
	for _, v := range list {
		uints = append(uints, v.Id)
	}
	return uints
}

func (a *AdminRolesGroup) DelByRolesGroupId(groupId uint) {
	Db().Where("roles_group_id = ?", groupId).Delete(&a)
}
