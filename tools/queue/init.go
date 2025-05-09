package queue

import (
	"app/internal/model"
	"app/tools/logger"
	"encoding/json"
	"errors"
)

func init() {
	registerHook()
	registerSaveUserCallback()
}

func registerSaveUserCallback() {
	// 注册回调
	var saveUserCallback CallbackFunc = func(msg *Msg) *CallbackResult {
		// 业务逻辑
		user := &model.User{}
		err := json.Unmarshal([]byte(msg.Data), user)
		if err != nil {
			return &CallbackResult{
				Err:      err,
				Msg:      "json unmarshal err",
				Code:     1,
				BackData: msg,
			}
		}

		model.Db().Select("nickname", "password", "age", "sex", "token", "status").Model(user).Updates(user)

		return &CallbackResult{
			Err:      nil,
			Msg:      "success",
			Code:     0, // 0 成功，1 失败
			BackData: nil,
		}
	}
	RegisterCallback("saveUser", &saveUserCallback)
}

// 注册钩子
func registerHook() {
	var u HookFunc = func(stream Stream, data map[string]any) *HookResult {
		_, ok := data["msg"]
		if !ok {
			return &HookResult{
				Err:      errors.New("nil msg"),
				Msg:      "nil msg",
				Code:     1,
				BackData: nil,
			}
		}
		msg := data["msg"].(*Msg)

		_, ok = data["consumer"]
		if !ok {
			return &HookResult{
				Err:      errors.New("nil consumer"),
				Msg:      "nil consumer",
				Code:     1,
				BackData: nil,
			}
		}
		consumer := data["consumer"].(string)
		logger.System("CALLBACK MSG", "Msg", msg.Id, "consumer", consumer)
		return &HookResult{
			Err:      nil,
			Msg:      "success",
			Code:     0,
			BackData: nil,
		}
	}
	RegisterHook(CallbackSuccess, &u)
}
