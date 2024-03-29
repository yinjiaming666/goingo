package model

import (
	"goingo/tools"
	"goingo/tools/resp"
	"time"
)

type User struct {
	*MysqlBaseModel `gorm:"-:all"` // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint           `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Password        string         `json:"password" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Nickname        string         `json:"nickname" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	HeadImg         string         `json:"head_img" gorm:"type:VARCHAR(255) NOT NULL;default:''"`
	Status          uint8          `json:"status" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	CreateTime      int64          `json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"` // 自动写入时间戳
	CreateTimeStr   string         `json:"create_time_str" gorm:"-:all"`                                     // -:all 无读写迁移权限，该字段不在数据库中
	Age             uint8          `json:"age" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	Sex             uint8          `json:"sex" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
}

func (user *User) GetUserInfo() *User {
	Db().Where(&user).Debug().First(&user)
	return user
}

func (user *User) DoRegister() uint {
	user.Password = tools.Md5(user.Password, UserPwdSalt)
	user.CreateTime = int64(uint(time.Now().Unix()))
	tx := Db().Create(&user)
	if tx.Error != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "错误", Body: map[string]any{
			"error": tx.Error.Error(),
		}}).Response()
	}
	return user.Id
}

func (user *User) SetUser() *User {
	if user.Password != "" {
		user.Password = tools.Md5(user.Password, UserPwdSalt)
	}
	Db().Select("nickname", "password", "age", "sex", "token", "status").Model(&user).Updates(&user)
	return user
}
