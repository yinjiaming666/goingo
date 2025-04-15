package logic

import (
	model2 "app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"context"
	"errors"
	"fmt"
	"time"
)

type TokenLogic struct {
}

var TokenLogicInstance *TokenLogic

func init() {
	TokenLogicInstance = &TokenLogic{}
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
	tokenModel.DelToken() // 删除这个用户的 token

	cacheKey := model2.KeyUtils.GetTokenKey(userJwt.Token)
	uidTokenKey := model2.KeyUtils.GetUidToken(int(uid))

	// 删除旧的 key
	get := model2.RedisClient.Get(context.Background(), uidTokenKey)
	if get.Val() != "" {
		model2.RedisClient.Del(context.Background(), get.Val())
	}

	tokenModel.CreateToken()

	m := conv.Struct2Map(*userJwt, true)
	m["type"], _ = conv.Conv[string](m["type"])
	_, err := model2.RedisClient.HMSet(context.Background(), cacheKey, m).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "jwt 缓存失败", Body: nil}).Response()
	}
	_, err = model2.RedisClient.Set(context.Background(), uidTokenKey, cacheKey, -1).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "uidTokenKey 缓存失败", Body: nil}).Response()
	}
	if exTime > 0 {
		model2.RedisClient.Expire(context.Background(), cacheKey, time.Duration(exTime)*time.Second)
	}
	return j, userJwt
}

func (tl *TokenLogic) CheckJwt(j string) (*jwt.UserJwt, error) {
	userJwt, err := jwt.ParseJwt(j)
	if err != nil {
		return nil, err
	}

	cacheKey := model2.KeyUtils.GetTokenKey(userJwt.Token)
	r, err := model2.RedisClient.HGetAll(context.Background(), cacheKey).Result()
	if err != nil {
		return nil, err
	}
	i, ok := r["uid"]
	if !ok {
		return nil, errors.New("账户已经在其他终端上登录")
	}

	uid, _ := conv.Conv[uint](i)
	fmt.Println(r)
	fmt.Println(userJwt)
	fmt.Println(uid)
	if uid != userJwt.Uid || r["token"] != userJwt.Token || r["type"] != string(userJwt.Type) {
		return nil, errors.New("账户已经在其他终端上登录[1]")
	}

	if userJwt.ExpireTime < time.Now().Unix() {
		(&resp.JsonResp{Code: resp.ReFail, Message: "token 过期", Body: nil}).Response()
	}
	return userJwt, nil
}
