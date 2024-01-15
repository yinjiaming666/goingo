package logic

import (
	"fmt"
	model2 "goingo/internal/model"
	"goingo/tools/conv"
	"goingo/tools/resp"
	"time"
)

type UserLogic struct {
}

// LoadUser 根据 uid 搜索用户
func (u UserLogic) LoadUser(uid uint) *model2.User {
	var cacheKey = model2.KeyUtils.GetUserKey(uid)
	result, err := model2.RedisClient.HGetAll(cacheKey).Result()
	if err != nil {
		resp.Resp(resp.ReFail, "查询失败", struct{ e error }{e: err})
	}

	strUid, ok := result["id"]
	if !ok || strUid == "" {
		user := model2.User{
			Id: uid,
		}
		u := user.GetUserInfo()
		go func() {
			model2.RedisClient.HMSet(cacheKey, u.ToMap(u))
			model2.RedisClient.Expire(cacheKey, 10*time.Second)
		}()
		return u
	}

	user := new(model2.User)
	user.InitWithMap(conv.Map2AnyMap[string](result), user)
	return user
}

// EditUserInfo todo
func (u UserLogic) EditUserInfo(user *model2.User) error {
	if user.Id == 0 {
		return fmt.Errorf("缺少 uid")
	}
	user = user.SetUser()
	var cacheKey = model2.KeyUtils.GetUserKey(user.Id)
	go func() {
		model2.RedisClient.HMSet(cacheKey, user.ToMap(user))
		model2.RedisClient.Expire(cacheKey, 10*time.Second)
	}()
	return nil
}

func (u UserLogic) SearchUser(search map[string]interface{}) *model2.User {
	user := model2.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	return user.GetUserInfo()
}
