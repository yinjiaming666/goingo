package logic

import (
	model2 "app/internal/model"
)

type UserLogic struct {
}

var UserLogicInstance UserLogic

func init() {
	UserLogicInstance = UserLogic{}
}

// LoadUser 根据 uid 搜索用户
func (u UserLogic) LoadUser(uid uint) *model2.User {
	userModel := &model2.User{Id: uid}
	userModel.GetCache()
	if userModel.Id == 0 {
		userModel.GetUserInfo()
	}
	return userModel
}

func (u UserLogic) SearchUser(search map[string]interface{}) *model2.User {
	user := &model2.User{}
	if _, ok := search["nickname"]; ok {
		v, ok := search["nickname"].(string)
		if ok {
			user.Nickname = v
		}
	}
	user.GetUserInfo()
	return user
}
