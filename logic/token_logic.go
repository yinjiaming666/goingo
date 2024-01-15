package logic

import (
	"fmt"
	"goingo/database"
	"goingo/utils/conv"
	"goingo/utils/jwt"
	"goingo/utils/resp"
	"time"
)

type TokenLogic struct {
}

func (tl *TokenLogic) GenerateJwt(uid uint, jType jwt.JType, exTime int64) (string, *jwt.UserJwt) {

	j, userJwt := jwt.CreateJwt(uid, jType, exTime)

	tokenModel := new(database.Token)
	tokenModel.Uid = userJwt.Uid
	tokenModel.Token = userJwt.Token
	tokenModel.ExpireTime = userJwt.ExpireTime
	tokenModel.Type = string(userJwt.Type)
	tokenModel.DeviceId = userJwt.DeviceId
	tokenModel.DeviceType = userJwt.DeviceType

	//tokenModel.DelToken() // 删除这个用户的 token
	cacheKey := database.KeyUtils.GetTokenKey(userJwt.Token)
	database.RedisClient.Del(cacheKey) // 删除旧的 key

	tokenModel.CreateToken()

	m := conv.Struct2Map(*userJwt, true)
	m["type"], _ = conv.Conv[string](m["type"])
	_, err := database.RedisClient.HMSet(cacheKey, m).Result()
	if err != nil {
		resp.Resp(resp.ReFail, "jwt 缓存失败", map[string]any{})
	}
	if exTime > 0 {
		database.RedisClient.Expire(cacheKey, time.Duration(exTime)*time.Second)
	}
	return j, userJwt
}

func (tl *TokenLogic) CheckJwt(j string) *jwt.UserJwt {
	userJwt, err := jwt.ParseJwt(j)
	if err != nil {
		resp.Resp(resp.ReFail, "jwt 解析失败", map[string]any{})
	}

	cacheKey := database.KeyUtils.GetTokenKey(userJwt.Token)
	r, err := database.RedisClient.HGetAll(cacheKey).Result()
	if err != nil {
		resp.Resp(resp.ReFail, "查询失败", map[string]any{})
	}
	i, ok := r["uid"]
	if !ok {
		resp.Resp(resp.ReFail, "非法的 jwt", map[string]any{})
	}

	uid, _ := conv.Conv[uint](i)
	fmt.Println(r)
	fmt.Println(userJwt)
	fmt.Println(uid)
	if uid != userJwt.Uid || r["token"] != userJwt.Token || r["type"] != string(userJwt.Type) {
		resp.Resp(resp.ReFail, "非法的 jwt [1]", map[string]any{})
	}

	if userJwt.ExpireTime < time.Now().Unix() {
		resp.Resp(resp.ReFail, "token 过期", nil)
	}
	return userJwt
}
