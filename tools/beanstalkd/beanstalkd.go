package beanstalkd

import (
	"app/tools/logger"
	"encoding/json"
	"errors"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

var MsgChannel chan *Msg

type Beanstalkd struct {
	Ip     string
	Port   string
	client *beanstalk.Conn
}

var Instance *Beanstalkd

func init() {
	Instance = &Beanstalkd{}
}

func (b *Beanstalkd) Init(ip, port string) error {
	if b.client != nil {
		return errors.New("beanstalkd already initialized")
	}

	var err error
	b.client, err = beanstalk.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	return nil
}

func (b *Beanstalkd) Push(data *Msg, delay time.Duration) (uint64, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	id, err := b.client.Put(marshal, 1, delay, time.Minute)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (b *Beanstalkd) ReserveLoop() {
	for {
		id, body, err := b.client.Reserve(1 * time.Second)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		logger.Debug("beanstalkd reserve:", "id", id, "body", string(body))

		data := &Msg{
			JobId: id,
		}
		err = json.Unmarshal(body, data)

		MsgChannel <- data

		if err != nil {
			logger.Error("beanstalkd reserve Unmarshal err:", "err", err)
			return
		}
	}
}

func (b *Beanstalkd) HandelJob(id uint64, t string) {
	var err error
	if t == "delete" {
		err = b.client.Delete(id)
	}
	if err != nil {
		logger.Error("beanstalkd Delete err:", "err", err, "handel type", t, "handel id", id)
		return
	}
}

type C string
type M string

var CDefault C
var MSaveUser M

type Msg struct {
	C     C      `json:"c"`
	M     M      `json:"m"`
	Data  any    `json:"data"`
	JobId uint64 `json:"job_id"`
}
