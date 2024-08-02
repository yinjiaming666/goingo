package logic

import (
	model2 "app/internal/model"
	"app/tools/conv"
	"app/tools/resp"
	"context"
	"fmt"
	"time"
)

type UserLogic struct {
}

// LoadUser 根据 uid 搜索用户
func (u UserLogic) LoadUser(uid uint) *model2.User {
	var cacheKey = model2.KeyUtils.GetUserKey(uid)
	result, err := model2.RedisClient.HGetAll(context.Background(), cacheKey).Result()
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "查询失败", Body: map[string]any{"error": err.Error()}}).Response()
	}

	strUid, ok := result["id"]
	if !ok || strUid == "" {
		user := model2.User{
			Id: uid,
		}
		u := user.GetUserInfo()
		go func() {
			model2.RedisClient.HMSet(context.Background(), cacheKey, u.ToMap(u))
			model2.RedisClient.Expire(context.Background(), cacheKey, 10*time.Second)
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
		model2.RedisClient.HMSet(context.Background(), cacheKey, user.ToMap(user))
		model2.RedisClient.Expire(context.Background(), cacheKey, 10*time.Second)
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
