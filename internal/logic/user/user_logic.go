package user

import (
	model2 "app/internal/model"
)

// LoadUser 根据 uid 搜索用户
func LoadUser(uid uint) *model2.User {
	userModel := &model2.User{Id: uid}
	userModel.GetCache()
	if userModel.Id == 0 {
		userModel.Id = uid
		userModel.GetUserInfo(true)
	}
	return userModel
}

func SearchUser(search map[string]interface{}) *model2.User {
	user := &model2.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	user.GetUserInfo(false)
	return user
}
