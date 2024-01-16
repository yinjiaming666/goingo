package model

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goingo/tools/key_utils"
	"goingo/tools/logger"
)

// RedisClient 定义一个全局变量
var RedisClient = &redis.Client{}
var KeyUtils = &key_utils.KeyUtils{}

type RedisConf struct {
	Ip         string
	Port       string
	GlobalName string
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

	KeyUtils.BaseName = c.GlobalName
}
