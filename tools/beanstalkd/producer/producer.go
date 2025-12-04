package producer

import (
	"app/tools/beanstalkd/message"
	"encoding/json"
	"errors"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// BeanstalkdProducer 生产者
type BeanstalkdProducer struct {
	Ip       string
	Port     string
	client   *beanstalk.Conn
	TubeName string
	Tube     *beanstalk.Tube
}

var Instance *BeanstalkdProducer

func init() {
	Instance = &BeanstalkdProducer{}
}

func (b *BeanstalkdProducer) Init(ip, port, TubeName string) error {
	if b.client != nil {
		return errors.New("beanstalkd producer already initialized")
	}
	var err error
	b.client, err = beanstalk.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	b.TubeName = TubeName
	if b.TubeName != "" {
		b.Tube = &beanstalk.Tube{Conn: b.client, Name: TubeName}
	}
	return nil
}

func (b *BeanstalkdProducer) Push(data *message.Message, delay time.Duration) (uint64, error) {
	if !json.Valid([]byte(data.Data)) {
		return 0, errors.New("beanstalkd push is not json data [" + data.Data + "]")
	}
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

func (b *BeanstalkdProducer) PushTube(data *message.Message, delay time.Duration) (uint64, error) {
	if !json.Valid([]byte(data.Data)) {
		return 0, errors.New("beanstalkd push tube is not json data [" + data.Data + "]")
	}
	marshal, err := json.Marshal(data)
	id, err := b.Tube.Put(marshal, 1, delay, time.Minute)
	if err != nil {
		return 0, err
	}
	return id, nil
}
