package main

import (
	global "app/internal"
	"app/internal/model"
	"app/internal/router"
	confg "app/tools/config"
	"app/tools/logger"
	"app/tools/queue"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const LOGO = `
 ██████╗   ██████╗  ██╗ ███╗   ██╗  ██████╗   ██████╗ 
██╔════╝  ██╔═══██╗ ██║ ████╗  ██║ ██╔════╝  ██╔═══██╗
██║  ███╗ ██║   ██║ ██║ ██╔██╗ ██║ ██║  ███╗ ██║   ██║
██║   ██║ ██║   ██║ ██║ ██║╚██╗██║ ██║   ██║ ██║   ██║
╚██████╔╝ ╚██████╔╝ ██║ ██║ ╚████║ ╚██████╔╝ ╚██████╔╝
 ╚═════╝   ╚═════╝  ╚═╝ ╚═╝  ╚═══╝  ╚═════╝   ╚═════╝ 
`

var err error

func main() {
	fmt.Print(LOGO)

	flag.StringVar(&global.Mode, "mode", "dev", "-mode=prod, -mode=dev") // "dev" or "prod"
	flag.StringVar(&global.InitDb, "initDb", "false", "-initDb=true, -initDb=false")
	flag.Parse()
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	conf := (&confg.Config{
		Path:     "./config",
		FileName: global.Mode, // dev or prod
	}).Init()

	global.ServerName = confg.Get[string](conf, "server", "name")
	pid := os.Getpid()
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(pid))
	f, _ := os.Create(global.ServerName + ".pid")
	_, err = f.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Printf("进程 PID: %d 写入失败 \n", pid)
		return
	}
	fmt.Printf("进程 PID: %d \n", pid)

	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本地 ip 失败" + err.Error())
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrList {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				global.LocalIp = ipNet.IP.String()
				break
			}
		}
	}

	logger.Init()

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

	if global.InitDb == "true" {
		logger.Info("start init table ====================")
		m := new(model.MysqlBaseModel)
		m.CreateTable(model.User{})
		m.CreateTable(model.Token{})
		m.CreateTable(model.Article{})
		m.CreateTable(model.Admin{})
		m.CreateTable(model.Cate{})
		m.CreateTable(model.Menu{})
		m.CreateTable(model.Role{})
		logger.Info("end init table ====================")
	}

	//queue.Init("goingo-queue", redis.NewClient(&redis.Options{
	//	Addr: confg.Get[string](conf, "queue", "ip") + ":" + confg.Get[string](conf, "queue", "port"),
	//}))
	//
	//// 消息队列
	//stream := &queue.NormalStream{}
	//stream.SetName("default")
	//err = stream.Create()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//go stream.Loop()
	//
	//// 延时队列
	//delayStream := &queue.DelayStream{}
	//delayStream.SetName("default")
	//err = delayStream.Create()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//go delayStream.Loop()
	//
	//initQueueFunc()

	port := confg.Get[string](conf, "server", "port")
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
