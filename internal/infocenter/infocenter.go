package infocenter

import (
	"app/internal/model"
	"app/tools/beanstalkd"
	"app/tools/logger"
)

func init() {
	for {
		select {
		case body := <-beanstalkd.MsgChannel:
			switch body.M {
			case beanstalkd.MSaveUser:
				if user, ok := body.Data.(model.User); ok {
					model.Db().Select("nickname", "password", "age", "sex", "token", "status").Model(user).Updates(user)
					beanstalkd.Instance.HandelJob(body.JobId, "delete")
				} else {
					logger.Error("infocenter save user failed data type is not model.User")
				}
				break
			}
		}
	}
}
