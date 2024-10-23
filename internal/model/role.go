package model

import "strings"

// Role 角色表
type Role struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Pid             uint           `json:"pid" gorm:"type:INT(11) UNSIGNED NOT NULL;default:0"`
	MenuIds         string         `json:"menu_ids" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
}

func (r *Role) GetMenusByRoleIds(ids string, t int) []*MenusFormat {
	menuIds := ""
	if ids != "*" {
		list := make([]*Role, 0)
		Db().Model(r).Where("id in ?", strings.Split(ids, ",")).Find(&list)
		for _, role := range list {
			if menuIds == "" {
				menuIds = role.MenuIds
			} else {
				menuIds += "," + role.MenuIds
			}
		}
	} else {
		menuIds = "*"
	}
	if menuIds == "" {
		return make([]*MenusFormat, 0)
	}
	menuList := make([]*MenusFormat, 0)
	menu := Menu{}
	c := Db().Model(menu)
	if menuIds != "*" {
		c.Where("id in ?", strings.Split(menuIds, ","))
	}
	if t >= 0 {
		c.Where("type = ?", t)
	}
	c.Find(&menuList)
	return menu.FormatTree(menuList)
}
