package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goingo/internal/model"
	"goingo/internal/router"
	"goingo/tools"
	"goingo/tools/logger"
	"goingo/tools/queue"
	"time"
)

var mode = flag.String("mode", "dev", "-mode=prod,-mode=dev")

var serverName = tools.GetConfig(*mode, "server", "name")
var initDb = flag.String("initDb", "false", "-initDb=true, -initDb=false")

func main() {
	flag.Parse()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	logger.InitLog()
	model.InitDb(&model.DbConf{
		UserName: tools.GetConfig(*mode, "mysql", "username"),
		Password: tools.GetConfig(*mode, "mysql", "password"),
		Ip:       tools.GetConfig(*mode, "mysql", "ip"),
		Port:     tools.GetConfig(*mode, "mysql", "port"),
		DbName:   tools.GetConfig(*mode, "mysql", "db_name"),
	})

	model.InitRedis(&model.RedisConf{
		Ip:         tools.GetConfig(*mode, "redis", "ip"),
		Port:       tools.GetConfig(*mode, "redis", "port"),
		GlobalName: serverName,
	})

	if *initDb == "true" {
		logger.Info("start init table ====================")
		m := new(model.MysqlBaseModel)
		m.CreateTable(model.User{})
		m.CreateTable(model.Token{})
		m.CreateTable(model.Article{})
		m.CreateTable(model.Admin{})
		m.CreateTable(model.Cate{})
		logger.Info("end init table ====================")
	}

	queue.Init("goingo-queue", redis.NewClient(&redis.Options{
		Addr: tools.GetConfig(*mode, "queue", "ip") + ":" + tools.GetConfig(*mode, "queue", "port"),
	}))

	// 消息队列
	stream := &queue.NormalStream{}
	stream.SetName("default")
	err := stream.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go stream.Loop()

	// 延时队列
	delayStream := &queue.DelayStream{}
	delayStream.SetName("default")
	err = delayStream.Create()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go delayStream.Loop()

	initQueueFunc()

	port := tools.GetConfig(*mode, "server", "port")
	router.InitRouter(port)
}

func initQueueFunc() {
	// 注册回调
	var pF queue.CallbackFunc = func(msg *queue.Msg) *queue.CallbackResult {
		// 业务逻辑
		fmt.Println("callback")
		fmt.Println(msg)
		return &queue.CallbackResult{
			Err:      nil,
			Msg:      "",
			Code:     0, // 0 成功，1 失败
			BackData: nil,
		}
	}
	queue.RegisterCallback("echo", &pF)

	// 注册钩子
	var u queue.HookFunc = func(stream queue.Stream, data map[string]any) *queue.HookResult {
		_, ok := data["msg"]
		if !ok {
			return &queue.HookResult{
				Err:      errors.New("nil msg"),
				Msg:      "nil msg",
				Code:     1,
				BackData: nil,
			}
		}
		msg := data["msg"].(*queue.Msg)

		_, ok = data["consumer"]
		if !ok {
			return &queue.HookResult{
				Err:      errors.New("nil consumer"),
				Msg:      "nil consumer",
				Code:     1,
				BackData: nil,
			}
		}
		consumer := data["consumer"].(string)
		logger.System("CALLBACK MSG", "Msg", msg.Id, "consumer", consumer)
		return &queue.HookResult{
			Err:      nil,
			Msg:      "success",
			Code:     0,
			BackData: nil,
		}
	}
	queue.RegisterHook(queue.CallbackSuccess, &u)
}
