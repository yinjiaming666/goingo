package main

import (
	"flag"
	"goingo/database"
	"goingo/logger"
	"goingo/router"
	"goingo/utils"
	"time"
)

var mode = flag.String("mode", "dev", "-mode=prod,-mode=dev")

// var serverName = utils.GetConfig(*mode, "server", "name")
var initDb = flag.String("initDb", "false", "-initDb=true, -initDb=false")

func main() {
	flag.Parse()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc

	logger.InitLog()

	database.InitDb(&database.DbConf{
		UserName: utils.GetConfig(*mode, "mysql", "username"),
		Password: utils.GetConfig(*mode, "mysql", "password"),
		Ip:       utils.GetConfig(*mode, "mysql", "ip"),
		Port:     utils.GetConfig(*mode, "mysql", "port"),
		DbName:   utils.GetConfig(*mode, "mysql", "db_name"),
	})

	database.InitRedis(&database.RedisConf{
		Ip:         utils.GetConfig(*mode, "redis", "ip"),
		Port:       utils.GetConfig(*mode, "redis", "port"),
		GlobalName: utils.GetConfig(*mode, "server", "name"),
	})

	if *initDb == "true" {
		logger.Info("start init table ====================")
		model := new(database.MysqlBaseModel)
		model.CreateTable(database.User{})
		model.CreateTable(database.Token{})
		model.CreateTable(database.Article{})
		model.CreateTable(database.Admin{})
		model.CreateTable(database.Cate{})
		logger.Info("end init table ====================")
	}

	port := utils.GetConfig(*mode, "server", "port")
	router.InitRouter(port)
}
