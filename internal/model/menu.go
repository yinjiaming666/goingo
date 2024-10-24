package model

import (
	"database/sql/driver"
	"encoding/json"
	"sort"
)

// Menu 菜单表
type Menu struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Component       string         `json:"component" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Meta            MenuMeta       `json:"meta" gorm:"type:json;default:NULL"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Path            string         `json:"path" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Pid             uint           `json:"pid" gorm:"type:INT(11) UNSIGNED NOT NULL;default:0"`
	Type            uint           `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"` // 是否为菜单
}

type MenusFormat struct {
	Id        uint           `json:"id"`
	Component string         `json:"component"`
	Meta      MenuMeta       `json:"meta"`
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Pid       uint           `json:"pid"`
	Type      uint           `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"` // 是否为菜单
	Children  []*MenusFormat `json:"children,omitempty" gorm:"-"`
}

func (m *Menu) FormatTree(list []*MenusFormat) []*MenusFormat {
	tempMap := make(map[uint]*MenusFormat)
	tempList := make([]*MenusFormat, 0)
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

func (m *Menu) SetMenu() *Menu {
	if m.Id <= 0 {
		Db().Create(&m)
	} else {
		Db().Model(&m).Updates(&m)
	}
	return m
}

func (m *Menu) SearchByPath(path string) {
	Db().Model(m).Where("path = ?", path).Where("type = 0").First(m)
}

func (m *Menu) DelMenu(ids []int) {
	Db().Delete(&m, ids)
}

type MenuMeta struct {
	Title    string `json:"title,omitempty"`
	AffixTab bool   `json:"affix_tab,omitempty"` // 是否固定标签页
	Order    int    `json:"order,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

func (c MenuMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *MenuMeta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

func (m *Menu) GetMenus() []*MenusFormat {
	list := make([]*MenusFormat, 0)
	Db().Model(m).Where("pid = ?", 0).Find(&list)
	for _, v := range list {
		Db().Model(m).Where("pid = ?", v.Id).Find(&v.Children)
	}
	return list
}
