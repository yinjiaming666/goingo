package model

import (
	"database/sql/driver"
	"encoding/json"
	"sort"
)

// Roles 菜单表
type Roles struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey; type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Component       string         `json:"component" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Meta            RolesMeta      `json:"meta" gorm:"type:json; default:NULL"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Path            string         `json:"path" gorm:"type:VARCHAR(1200) NOT NULL; default:''"`
	Pid             uint           `json:"pid" gorm:"type:INT(11) UNSIGNED NOT NULL; default:0"`
	Type            uint           `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL; default:0"` // 是否为菜单
}

type RolesFormat struct {
	Id        uint           `json:"id"`
	Component string         `json:"component"`
	Meta      RolesMeta      `json:"meta"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Pid       uint           `json:"pid"`
	Type      uint           `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL; default:0; comment:是否为菜单"` // 是否为菜单
	Children  []*RolesFormat `json:"children,omitempty" gorm:"-"`
}

func (m *Roles) FormatTree(list []*RolesFormat) []*RolesFormat {
	tempMap := make(map[uint]*RolesFormat)
	tempList := make([]*RolesFormat, 0)
	for _, v := range list {
		tempMap[v.Id] = v
	}
	for _, v := range tempMap {
		if _, ok := tempMap[v.Pid]; ok {
			tempMap[v.Pid].Children = append(tempMap[v.Pid].Children, tempMap[v.Id])
		} else {
			tempList = append(tempList, tempMap[v.Id])
		}
	}

	sort.Slice(tempList, func(i, j int) bool { return tempList[i].Id < tempList[j].Id })
	return tempList
}

func (m *Roles) SetRoles() *Roles {
	if m.Id <= 0 {
		Db().Create(&m)
	} else {
		Db().Model(&m).Updates(&m)
	}
	return m
}

func (m *Roles) SearchByPath(path string) {
	Db().Model(m).Where("path = ?", path).Where("type = 0").First(m)
}

func (m *Roles) DelRoles(ids []int) {
	Db().Delete(&m, ids)
}

type RolesMeta struct {
	Title    string `json:"title,omitempty"`
	AffixTab bool   `json:"affix_tab,omitempty"` // 是否固定标签页
	Order    int    `json:"order,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

func (c RolesMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *RolesMeta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

func (m *Roles) GetRoles(ids []uint, t int) []*RolesFormat {
	tx := Db().Model(m)
	list := make([]*RolesFormat, 0)
	if len(ids) > 0 {
		tx.Where("id IN ?", ids)
	}
	if t >= 0 {
		tx.Where("type = ?", t)
	}
	tx.Find(&list)
	return list
}
