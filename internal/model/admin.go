package model

import "app/tools/conv"

type Admin struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Account         string         `json:"account" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Password        string         `json:"password" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Avatar          string         `json:"avatar" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	RoleIds         string         `json:"role_ids" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	RoleList        []*MenusFormat `json:"role_list" gorm:"-:all"`
}

func (admin *Admin) GetAdmin() *Admin {
	Db().Where(admin).First(admin)
	return admin
}

func (admin *Admin) GetList(pid uint) []*Admin {
	list := make([]*Admin, 0)
	menu := Menu{}
	Db().Model(admin).Find(&list)
	for _, v := range list {
		explode, _ := conv.Explode[uint](",", v.RoleIds)
		v.RoleList = menu.GetMenus(explode)
		v.RoleList = menu.FormatTree(v.RoleList)
	}
	return list
}

func (admin *Admin) UpdateAdmin(data map[string]interface{}) *Admin {
	Db().Model(&admin).Updates(data)
	return admin
}
