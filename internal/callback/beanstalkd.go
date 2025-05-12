package callback

import (
	"app/internal/model"
	"app/tools/beanstalkd/consumer"
	beanstalkdMsg "app/tools/beanstalkd/message"
	"app/tools/logger"
	"encoding/json"
)

func init() {
	consumer.Instance.SetCallback(callback)
}

func callback(message *beanstalkdMsg.Message) {
	switch message.M {
	case beanstalkdMsg.MSaveUser:
		dataStr := message.Data
		user := &model.User{}
		err := json.Unmarshal([]byte(dataStr), user)
		if err != nil {
			logger.Error("beanstalkd callback [save user] json decode fail" + err.Error())
		}
		model.Db().Select("nickname", "password", "age", "sex", "token", "status").Model(user).Create(user)
		consumer.Instance.HandelJob(message.JobId, "delete")
		break
	}
}
