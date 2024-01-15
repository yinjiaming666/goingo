package database

import (
	"github.com/go-redis/redis"
	"goingo/logger"
	"goingo/utils/key_utils"
)

// RedisClient 定义一个全局变量
var RedisClient = &redis.Client{}
var OpenRedis = true // 是否启用 redis
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
	res, err := RedisClient.Ping().Result()
	if err != nil {
		println(err.Error())
		logger.Error("redis connect fail", err.Error())
	} else {
		logger.Error("redis init success", res)
	}
	OpenRedis = true
	KeyUtils.BaseName = c.GlobalName
}
