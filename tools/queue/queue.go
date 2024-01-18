package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goingo/tools/logger"
	"time"
)

var Client *redis.Client
var GlobalName string
var streamList = make(map[string]*Stream)

type SType string

// Normal 普通队列
var Normal SType = "normal"

// Delay 延时队列
var Delay SType = "delay"

// Stream 流
type Stream struct {
	Name        string
	fullName    string // redis 里存的名字
	Type        SType
	BackGroup   *XGroup // 用来备份的消费者组
	HandelGroup *XGroup // 用户执行的消费者组
}

// XGroup 消费组
type XGroup struct {
	streamName   string
	name         string
	start        string
	ConsumerList []*XConsumer
}

// XConsumer 消费者
type XConsumer struct {
	Name       string
	groupName  string
	streamName string
}

func Init(name string, client *redis.Client) {
	Client = client
	GlobalName = name
}

func getFullStreamName(name string, sType SType) string {
	return GlobalName + ":" + string(sType) + ":" + name
}

func CreateStream(stream *Stream) error {
	if _, ok := streamList[stream.Name]; ok {
		return errors.New("repeat stream")
	}
	stream.fullName = getFullStreamName(stream.Name, stream.Type)

	if len(stream.HandelGroup.ConsumerList) == 0 {
		return errors.New("empty HandelGroup")
	}
	stream.HandelGroup.streamName = stream.fullName
	stream.HandelGroup.name = "handel"
	stream.HandelGroup.start = "$"
	Client.XGroupDestroy(context.Background(), stream.fullName, "")
	Client.XGroupDestroy(context.Background(), stream.fullName, "handel")
	Client.XGroupDestroy(context.Background(), stream.fullName, "back")

	res, err := Client.XGroupCreate(context.Background(), stream.fullName, stream.HandelGroup.name, stream.HandelGroup.start).Result()
	if err != nil {
		// todo
	}
	logger.Info("队列：" + stream.fullName + "创建执行消费者组" + res)
	for k, consumer := range stream.HandelGroup.ConsumerList {
		if consumer.Name == "" {
			return errors.New("empty consumer name")
		}
		stream.HandelGroup.ConsumerList[k].streamName = stream.fullName
		stream.HandelGroup.ConsumerList[k].groupName = stream.HandelGroup.name
		_, err = Client.XGroupCreateConsumer(context.Background(), stream.fullName, stream.HandelGroup.name, consumer.Name).Result()
		if err != nil {
			// todo
		}
		logger.Info("队列：" + stream.fullName + "执行消费者组创建消费者：" + consumer.Name)
	}

	if len(stream.BackGroup.ConsumerList) > 0 {
		stream.BackGroup.streamName = stream.fullName
		stream.BackGroup.name = "back"
		stream.BackGroup.start = "$"
		res, err := Client.XGroupCreate(context.Background(), stream.fullName, stream.BackGroup.name, stream.BackGroup.start).Result()
		if err != nil {
			// todo
		}
		logger.Info("队列：" + stream.fullName + "创建备份消费者组" + res)
		for k, consumer := range stream.BackGroup.ConsumerList {
			if consumer.Name == "" {
				return errors.New("empty consumer name")
			}
			stream.BackGroup.ConsumerList[k].streamName = stream.fullName
			stream.BackGroup.ConsumerList[k].groupName = stream.BackGroup.name
			_, err = Client.XGroupCreateConsumer(context.Background(), stream.fullName, stream.BackGroup.name, consumer.Name).Result()
			if err != nil {
				// todo
			}
			logger.Info("队列：" + stream.fullName + "备份消费者组创建消费者：" + consumer.Name)
		}
	}

	streamList[stream.Name] = stream
	stream.EchoInfo()
	return nil
}

type MsgBody struct {
}

// EchoInfo 输出 stream 全部信息
func (s *Stream) EchoInfo() {
	stream, err := Client.XInfoStream(context.Background(), s.fullName).Result()
	if err != nil {
		fmt.Printf("队列 %s 获取失败：%v \n", s.fullName, err)
		return
	}
	fmt.Printf(">队列名称：%s \n", s.fullName)
	fmt.Printf(">队列长度：%d \n", stream.Length)
	fmt.Printf(">FirstEntry：%v \n", stream.FirstEntry)
	fmt.Printf(">LastEntry：%v \n", stream.LastEntry)
	fmt.Printf(">RecordedFirstEntryID：%s \n", stream.RecordedFirstEntryID)

	groups, err := Client.XInfoGroups(context.Background(), s.fullName).Result()
	if err != nil {
		fmt.Printf("消费组获取失败：%v \n", err)
		return
	}
	for _, group := range groups {
		fmt.Printf(">    消费组名称：%s \n", group.Name)
		fmt.Printf(">    Consumers：%d \n", group.Consumers)
		fmt.Printf(">    Pending：%d \n", group.Pending)
		fmt.Printf(">    LastDeliveredID：%s \n", group.LastDeliveredID)
		fmt.Printf(">    EntriesRead：%d \n", group.EntriesRead)
		fmt.Printf(">    Lag：%d \n", group.Lag)
		consumers, _ := Client.XInfoConsumers(context.Background(), s.fullName, group.Name).Result()
		for _, consumer := range consumers {
			fmt.Printf(">        消费者名称：%s \n", consumer.Name)
			fmt.Printf(">        Pending：%d \n", consumer.Pending)
			fmt.Printf(">        Idle：%s \n", consumer.Idle)
			fmt.Printf(">        Inactive：%s \n", consumer.Inactive)
		}
	}
}

func GetPending(name, gType string) (*redis.XPending, error) {
	stream, ok := streamList[name]
	if !ok {
		return nil, errors.New("not found stream")
	}
	var group *XGroup
	if gType == "back" {
		group = stream.HandelGroup
	} else {
		group = stream.BackGroup
	}
	result, err := Client.XPending(context.Background(), stream.Name, group.name).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Push(body map[string]interface{}, queueName string, sType SType) (string, error) {
	var b = &redis.XAddArgs{
		Stream: getFullStreamName(queueName, sType),
		MaxLen: 0,
		ID:     "",
		Values: body,
	}
	return Client.XAdd(context.Background(), b).Result()
}

func (c *XConsumer) listen() {
	for {
		fmt.Printf("============== %s ==============\n阻塞读取中\n%s\n%s\n============== %s ==============\n\n\n", c.Name, c.streamName, c.groupName, c.Name)
		result, err := Client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    c.groupName,
			Consumer: c.Name,
			Streams:  []string{c.streamName, ">"},
			Count:    1,
			Block:    1,
		}).Result()
		if err != nil {
			//println(err)
		} else {
			fmt.Printf("============== %s ============\n读取结果 %+v \n============== %s ============\n\n\n", c.Name, result, c.Name)
		}
		time.Sleep(1 * time.Second)
	}
}

func Loop() {
	for _, stream := range streamList {
		for _, hc := range stream.HandelGroup.ConsumerList {
			go hc.listen()
		}
		for _, bc := range stream.BackGroup.ConsumerList {
			go bc.listen()
		}
	}
}
