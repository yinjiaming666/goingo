package model

import "goingo/tools/jwt"

type Token struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Uid             uint           `json:"uid" gorm:"type:INT(8) UNSIGNED NOT NULL;default:0"`
	ExpireTime      int64          `json:"expire_time" gorm:"type:BIGINT UNSIGNED NOT NULL;default:0"`
	Token           string         `json:"token" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Type            string         `json:"type" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	DeviceType      uint8          `json:"device_type" gorm:"type:TINYINT(8) NOT NULL;default:0"`  // 设备类型
	DeviceId        string         `json:"device_id" gorm:"type:VARCHAR(255) NOT NULL;default:''"` // 设备id
	CreateTime      int64          `json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"`
	CreateTimeStr   string         `json:"create_time_str" gorm:"-:all"`
}

func (t *Token) CheckToken(token string, jType jwt.JType) *Token {
	Db().First(t, "token = ? AND type = ?", token, string(jType))
	return t
}

func (t *Token) CreateToken() *Token {
	Db().Create(&t)
	return t
}

func (t *Token) DelToken() {
	where := make(map[string]interface{})
	if t.Uid > 0 {
		where["uid"] = t.Uid
	}
	if t.Id > 0 {
		where["id"] = t.Id
	}
	if t.Type != "" {
		where["type"] = t.Type
	}
	Db().Where(where).Delete(&t)
}
