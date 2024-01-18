package queue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goingo/tools/logger"
	"time"
)

var Client *redis.Client

type Queue struct {
	Client     *redis.Client
	GlobalName string
	streamList []Stream
}

// Stream 流
type Stream struct {
	Name      string
	GroupList []XGroup
}

// XGroup 消费组
type XGroup struct {
	StreamName   string
	Name         string
	Start        string
	ConsumerList []XConsumer
}

// XConsumer 消费者
type XConsumer struct {
	Name string
}

type MsgBody struct {
}

func (q *Queue) Init() *Queue {
	q.GlobalName += ":queue:"
	Client = q.Client
	q.SetStream(&Stream{
		Name: q.GlobalName + "default",
		GroupList: []XGroup{
			{
				StreamName: q.GlobalName + "default",
				Name:       "default_group_1",
				Start:      "$",
				ConsumerList: []XConsumer{
					{
						Name: "default_group_1_consumer_1",
					},
					{
						Name: "default_group_1_consumer_2",
					},
				},
			},
			{
				StreamName: q.GlobalName + "default",
				Name:       "default_group_2",
				Start:      "$",
				ConsumerList: []XConsumer{
					{
						Name: "default_group_2_consumer_1",
					},
					{
						Name: "default_group_2_consumer_2",
					},
				},
			},
		},
	})
	return q
}

func (q *Queue) SetStream(stream *Stream) {
	for _, group := range stream.GroupList {
		if group.Start == "" {
			group.Start = "$"
		}
		res, err := q.Client.XGroupCreate(context.Background(), stream.Name, group.Name, group.Start).Result()
		if err != nil {
			// todo
			//return
		}
		logger.Info("队列创建消费组", res)
		for _, consumer := range group.ConsumerList {
			res, err := q.Client.XGroupCreateConsumer(context.Background(), stream.Name, group.Name, consumer.Name).Result()
			if err != nil {
				// todo
				//return
			}
			logger.Info("队列创建消费者", res)
		}
	}

	if len(q.streamList) == 0 {
		q.streamList = append(q.streamList, *stream)
	}
	for k, s := range q.streamList {
		var isUpdate bool
		if s.Name == stream.Name {
			isUpdate = true
		}
		if isUpdate {
			q.streamList[k] = *stream
		} else {
			q.streamList = append(q.streamList, *stream)
		}
	}
	stream.EchoInfo()
}

// EchoInfo 输出 stream 全部信息
func (s *Stream) EchoInfo() {
	stream, err := Client.XInfoStream(context.Background(), s.Name).Result()
	if err != nil {
		fmt.Printf("队列%s获取失败：%v \n", s.Name, err)
		return
	}
	fmt.Printf(">队列名称：%s \n", s.Name)
	fmt.Printf(">队列长度：%d \n", stream.Length)
	fmt.Printf(">FirstEntry：%v \n", stream.FirstEntry)
	fmt.Printf(">LastEntry：%v \n", stream.LastEntry)
	fmt.Printf(">RecordedFirstEntryID：%s \n", stream.RecordedFirstEntryID)

	groups, err := Client.XInfoGroups(context.Background(), s.Name).Result()
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
		consumers, _ := Client.XInfoConsumers(context.Background(), s.Name, group.Name).Result()
		for _, consumer := range consumers {
			fmt.Printf(">        消费者名称：%s \n", consumer.Name)
			fmt.Printf(">        Pending：%d \n", consumer.Pending)
			fmt.Printf(">        Idle：%s \n", consumer.Idle)
			fmt.Printf(">        Inactive：%s \n", consumer.Inactive)
		}
	}

}

func (q *Queue) GetPending(name string) *map[string]map[string]int64 {
	var list map[string]map[string]int64
	var stream Stream
	for _, s := range q.streamList {
		if s.Name == name {
			stream = s
		}
	}

	for _, group := range stream.GroupList {
		result, err := q.Client.XPending(context.Background(), stream.Name, group.Name).Result()
		if err != nil {

		}
		list[group.Name] = result.Consumers
	}
	return &list
}

func (q *Queue) Push(body map[string]interface{}, queueName string) (string, error) {
	var b = &redis.XAddArgs{
		Stream: q.GlobalName + queueName,
		MaxLen: 0,
		ID:     "",
		Values: body,
	}
	return q.Client.XAdd(context.Background(), b).Result()
}

func (g *XGroup) listen() {
	for {
		fmt.Println("du ====" + g.Name)
		result, err := Client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    g.Name,
			Consumer: g.ConsumerList[0].Name,
			Streams:  []string{g.StreamName, ">"},
			Count:    1,
		}).Result()
		if err != nil {
			fmt.Println("队列读取失败", err)
		}
		fmt.Printf("消费组 %s 读取结果：%+v \n", g.Name, result)
		time.Sleep(1 * 500)
	}
}

func (q *Queue) Loop() {
	for _, s := range q.streamList {
		for _, group := range s.GroupList {
			fmt.Println(group.Name)
			go group.listen()
		}
	}
}
