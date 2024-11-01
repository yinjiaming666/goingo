package jwt

import (
	sysLog "app/tools/logger"
	"app/tools/random"
	"github.com/golang-jwt/jwt"
	"time"
)

type JType string

const (
	AdminJwtType   JType = "admin_template"
	IndexJwtType   JType = "index_template"
	Key                  = "admin123456!"
	DefaultExpTime       = 30 * 86400 // jwt 默认过期时间（秒）
)

type UserJwt struct {
	Uid                uint   `json:"uid"`
	Type               JType  `json:"type"`
	Token              string `json:"token"`
	ExpireTime         int64  `json:"expire_time"`
	DeviceId           string `json:"device_id"`
	DeviceType         uint8  `json:"device_type"`
	jwt.StandardClaims        // 必须要实现这个接口
}

// CreateJwt 生成 jwt
func CreateJwt(id uint, jwtType JType, expireTime int64) (string, *UserJwt) {
	if expireTime < 0 {
		expireTime = time.Now().Unix() + (60 * 60 * 24 * 30)
	} else if expireTime == 0 {
		expireTime = time.Now().Unix() + DefaultExpTime
	} else {
		expireTime = time.Now().Unix() + expireTime
	}

	userJwt := UserJwt{
		Uid:            id,
		Type:           jwtType,
		Token:          random.Str(0),
		ExpireTime:     expireTime,
		DeviceId:       "",
		DeviceType:     0,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expireTime},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userJwt)
	tokenString, err := token.SignedString([]byte(Key))
	if err != nil {
		sysLog.Error("jwt 生成失败", err.Error())
		panic("jwt 生成失败" + err.Error())
	}
	return tokenString, &userJwt
}

func ParseJwt(token string) (*UserJwt, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserJwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Key), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserJwt); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
