// 队列消费者进程
package main

import (
	"app/internal/global"
	"app/internal/model"
	"app/tools/beanstalkd/callback"
	"app/tools/beanstalkd/consumer"
	"app/tools/beanstalkd/message"
	confg "app/tools/config"
	"app/tools/logger"
	"encoding/json"
	"fmt"
)

func main() {
	conf := (&confg.Config{
		Path:     "./config",
		FileName: global.Mode, // dev or prod
	}).Init()

	model.InitDb(&model.DbConf{
		UserName: confg.Get[string](conf, "mysql", "username"),
		Password: confg.Get[string](conf, "mysql", "password"),
		Ip:       confg.Get[string](conf, "mysql", "ip"),
		Port:     confg.Get[string](conf, "mysql", "port"),
		DbName:   confg.Get[string](conf, "mysql", "db_name"),
	})

	model.InitRedis(&model.RedisConf{
		Ip:   confg.Get[string](conf, "redis", "ip"),
		Port: confg.Get[string](conf, "redis", "port"),
	})

	beanstalkdIp := confg.Get[string](conf, "beanstalkd", "ip")
	beanstalkdPort := confg.Get[string](conf, "beanstalkd", "port")
	consumerNumber := confg.Get[int](conf, "beanstalkd", "consumer")

	if beanstalkdIp != "" && beanstalkdPort != "" {
		err := consumer.Instance.Init(beanstalkdIp, beanstalkdPort, []string{"common"})
		if err != nil {
			logger.Error("beanstalkd consumer init err:", "err", err)
			return
		}

		consumer.Instance.Callback = func(msg *message.Message) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("consumer ReserveLoop err", "err", err, "msg", fmt.Sprintf("%v", msg))
				}
			}()
			switch msg.M {
			case message.MHandelMoney:
				b := &message.HandelMoneyMsg{}
				err = json.Unmarshal([]byte(msg.Data), b)
				if err != nil {
					logger.Error("HandelMoneyMsg decode fail", "err", err, "data", msg.Data)
					break
				}
				callback.PutMoneyLog(b, msg.JobId)
				break
			}
		}

		if consumerNumber <= 0 {
			consumerNumber = 1
		}
		for i := 1; i <= consumerNumber; i++ {
			go func(i int) {
				consumer.Instance.ReserveLoop(consumerNumber)
			}(i)
		}
	}
	select {}
}
