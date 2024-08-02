package logic

import (
	model2 "app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"context"
	"fmt"
	"time"
)

type TokenLogic struct {
}

func (tl *TokenLogic) GenerateJwt(uid uint, jType jwt.JType, exTime int64) (string, *jwt.UserJwt) {

	j, userJwt := jwt.CreateJwt(uid, jType, exTime)

	tokenModel := new(model2.Token)
	tokenModel.Uid = userJwt.Uid
	tokenModel.Token = userJwt.Token
	tokenModel.ExpireTime = userJwt.ExpireTime
	tokenModel.Type = string(userJwt.Type)
	tokenModel.DeviceId = userJwt.DeviceId
	tokenModel.DeviceType = userJwt.DeviceType

	//tokenModel.DelToken() // 删除这个用户的 token
	cacheKey := model2.KeyUtils.GetTokenKey(userJwt.Token)
	model2.RedisClient.Del(context.Background(), cacheKey) // 删除旧的 key

	tokenModel.CreateToken()

	m := conv.Struct2Map(*userJwt, true)
	m["type"], _ = conv.Conv[string](m["type"])
	_, err := model2.RedisClient.HMSet(context.Background(), cacheKey, m).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "jwt 缓存失败", Body: nil}).Response()
	}
	if exTime > 0 {
		model2.RedisClient.Expire(context.Background(), cacheKey, time.Duration(exTime)*time.Second)
	}
	return j, userJwt
}

func (tl *TokenLogic) CheckJwt(j string) *jwt.UserJwt {
	userJwt, err := jwt.ParseJwt(j)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "jwt 解析失败", Body: nil}).Response()
	}

	cacheKey := model2.KeyUtils.GetTokenKey(userJwt.Token)
	r, err := model2.RedisClient.HGetAll(context.Background(), cacheKey).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "查询失败", Body: nil}).Response()
	}
	i, ok := r["uid"]
	if !ok {
		(&resp.JsonResp{Code: resp.ReFail, Message: "非法的 jwt", Body: nil}).Response()
	}

	uid, _ := conv.Conv[uint](i)
	fmt.Println(r)
	fmt.Println(userJwt)
	fmt.Println(uid)
	if uid != userJwt.Uid || r["token"] != userJwt.Token || r["type"] != string(userJwt.Type) {
		(&resp.JsonResp{Code: resp.ReFail, Message: "非法的 jwt [1]", Body: nil}).Response()
	}

	if userJwt.ExpireTime < time.Now().Unix() {
		(&resp.JsonResp{Code: resp.ReFail, Message: "token 过期", Body: nil}).Response()
	}
	return userJwt
}
