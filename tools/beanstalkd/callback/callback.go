package callback

import (
	"app/tools/beanstalkd/consumer"
	"app/tools/beanstalkd/message"
)

func PutMoneyLog(msg *message.HandelMoneyMsg, jobId uint64) {
	consumer.Instance.HandelJob(jobId, "delete")
}
