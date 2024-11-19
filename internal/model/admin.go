package model

type Admin struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey; type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Pid             uint           `json:"pid" gorm:"type:INT(11) UNSIGNED NOT NULL; default:0"`
	Account         string         `json:"account" gorm:"type:VARCHAR(1200) NOT NULL;  default:''; comment:登录账号"`
	Password        string         `json:"password" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Avatar          string         `json:"avatar" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	IsSuper         uint8          `json:"is_super" gorm:"type:TINYINT(8) NOT NULL; default:0; comment:1 超级管理员"`
	// RoleList        []*RolesFormat `json:"role_list" gorm:"-:all"`
}

func (admin *Admin) GetAdmin() *Admin {
	Db().Where(admin).First(admin)
	return admin
}

func (admin *Admin) GetList(pid uint) []*Admin {
	list := make([]*Admin, 0)
	//menu := Roles{}
	//Db().Model(admin).Find(&list)
	//for _, v := range list {
	//	explode, _ := conv.Explode[uint](",", v.RoleGroupIds)
	//	v.RoleList = menu.GetRoles(explode)
	//	v.RoleList = menu.FormatTree(v.RoleList)
	//}
	return list
}

func (admin *Admin) SetAdmin() *Admin {
	if admin.Id <= 0 {
		Db().Create(&admin)
	} else {
		Db().Select("title", "content", "status", "cate_id", "type").Model(&admin).Updates(&admin)
	}
	return admin
}

func (admin *Admin) DelAdmin(ids []int) {
	Db().Delete(&admin, ids)
}
