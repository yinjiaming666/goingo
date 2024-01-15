package model

type Admin struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Account         string         `json:"account" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Password        string         `json:"password" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Avatar          string         `json:"avatar" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Name            string         `json:"name" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
}

func (admin *Admin) GetAdmin() *Admin {
	db.Where(admin).Debug().First(admin)
	return admin
}

func (admin *Admin) UpdateAdmin(data map[string]interface{}) *Admin {
	db.Model(&admin).Debug().Updates(data)
	return admin
}
