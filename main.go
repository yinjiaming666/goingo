package main

import (
	"flag"
	"goingo/internal/model"
	"goingo/internal/router"
	"goingo/tools"
	"goingo/tools/logger"
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

	port := tools.GetConfig(*mode, "server", "port")
	router.InitRouter(port)
}
