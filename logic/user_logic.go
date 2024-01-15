package logic

import (
	"fmt"
	"goingo/database"
	"goingo/utils/conv"
	"goingo/utils/resp"
	"time"
)

type UserLogic struct {
}

// LoadUser 根据 uid 搜索用户
func (u UserLogic) LoadUser(uid uint) *database.User {
	var cacheKey = database.KeyUtils.GetUserKey(uid)
	result, err := database.RedisClient.HGetAll(cacheKey).Result()
	if err != nil {
		resp.Resp(resp.ReFail, "查询失败", struct{ e error }{e: err})
	}

	strUid, ok := result["id"]
	if !ok || strUid == "" {
		user := database.User{
			Id: uid,
		}
		u := user.GetUserInfo()
		go func() {
			database.RedisClient.HMSet(cacheKey, u.ToMap(u))
			database.RedisClient.Expire(cacheKey, 10*time.Second)
		}()
		return u
	}

	user := new(database.User)
	user.InitWithMap(conv.Map2AnyMap[string](result), user)
	return user
}

// EditUserInfo todo
func (u UserLogic) EditUserInfo(user *database.User) error {
	if user.Id == 0 {
		return fmt.Errorf("缺少 uid")
	}
	user = user.SetUser()
	var cacheKey = database.KeyUtils.GetUserKey(user.Id)
	go func() {
		database.RedisClient.HMSet(cacheKey, user.ToMap(user))
		database.RedisClient.Expire(cacheKey, 10*time.Second)
	}()
	return nil
}

func (u UserLogic) SearchUser(search map[string]interface{}) *database.User {
	user := database.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	return user.GetUserInfo()
}
