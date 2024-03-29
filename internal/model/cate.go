package model

import (
	"gorm.io/plugin/soft_delete"
)

type Cate struct {
	*MysqlBaseModel `gorm:"-:all"`        // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint                  `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Title           string                `json:"title"  gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Img             string                `json:"img"  gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	DeleteTime      soft_delete.DeletedAt `json:"delete_time" gorm:"type:BIGINT UNSIGNED NOT NULL;default:0"`
}

func (cate *Cate) GetCateList() []Cate {
	list := make([]Cate, 0)
	Db().Find(&list)
	return list
}
