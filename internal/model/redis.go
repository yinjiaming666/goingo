package model

import (
	global "app/internal"
	"app/tools/key_utils"
	"app/tools/logger"
	"context"
	"github.com/redis/go-redis/v9"
)

// RedisClient 定义一个全局变量
var RedisClient = &redis.Client{}
var KeyUtils = &key_utils.KeyUtils{}

type RedisConf struct {
	Ip   string
	Port string
}

func InitRedis(c *RedisConf) {
	o := &redis.Options{
		Addr: c.Ip + ":" + c.Port,
	}

	RedisClient = redis.NewClient(o)
	res, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		println(err.Error())
		logger.Error("redis connect fail", err.Error())
	} else {
		logger.Error("redis init success", res)
	}

	KeyUtils.BaseName = global.ServerName
}
