package consumer

import (
	"app/tools/beanstalkd/message"
	"app/tools/logger"
	"encoding/json"
	"errors"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

type CallbackFunc func(msg *message.Message)

// BeanstalkdConsumer 消费者
type BeanstalkdConsumer struct {
	Ip       string
	Port     string
	client   *beanstalk.Conn
	TubeName []string // 多个管道名称
	Callback CallbackFunc
}

var Instance *BeanstalkdConsumer

func init() {
	Instance = &BeanstalkdConsumer{}
}

func (b *BeanstalkdConsumer) Init(ip, port string, TubeName []string) error {
	if b.client != nil {
		return errors.New("beanstalkd consumer already initialized")
	}

	var err error
	b.client, err = beanstalk.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	b.TubeName = TubeName

	return nil
}

func (b *BeanstalkdConsumer) ReserveLoop() {
	if len(b.TubeName) > 0 {
		tubeSet := beanstalk.NewTubeSet(b.client, b.TubeName...)
		for {
			id, body, err := tubeSet.Reserve(1 * time.Second)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			logger.Debug("beanstalkd tube reserve:", "id", id, "body", string(body))

			data := &message.Message{}
			err = json.Unmarshal(body, data)
			data.JobId = id
			if err != nil {
				logger.Error("beanstalkd reserve Unmarshal err:", "err", err)
				continue
			} else {
				// 执行回调
				(b.Callback)(data)
			}
		}
	} else {
		for {
			id, body, err := b.client.Reserve(1 * time.Second)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			logger.Debug("beanstalkd tube reserve:", "id", id, "body", string(body))

			data := &message.Message{}
			err = json.Unmarshal(body, data)
			data.JobId = id
			if err != nil {
				logger.Error("beanstalkd tube reserve Unmarshal err:", "err", err)
				continue
			} else {
				// 执行回调
				(b.Callback)(data)
			}
		}
	}
}

func (b *BeanstalkdConsumer) HandelJob(id uint64, t string) {
	var err error
	if t == "delete" {
		err = b.client.Delete(id)
	}
	if err != nil {
		logger.Error("beanstalkd Delete err:", "err", err, "handel type", t, "handel id", id)
		return
	}
}

func (b *BeanstalkdConsumer) SetCallback(fun CallbackFunc) {
	b.Callback = fun
}
