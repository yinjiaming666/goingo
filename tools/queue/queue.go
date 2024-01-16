package queue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goingo/tools/logger"
)

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
	//q.DelStreamGroup("queue:default")
	q.GlobalName += ":queue:"
	q.SetStream(&Stream{
		Name: q.GlobalName + "default",
		GroupList: []XGroup{
			{
				Name:  "default_group_1",
				Start: "$",
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
				Name:  "default_group_2",
				Start: "$",
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
			fmt.Println("err1 ============")
			fmt.Println(stream.Name)
			fmt.Println(err)
			//return
		}
		logger.Info("队列创建消费组", res)
		for _, consumer := range group.ConsumerList {
			res, err := q.Client.XGroupCreateConsumer(context.Background(), stream.Name, group.Name, consumer.Name).Result()
			if err != nil {
				// todo
				fmt.Println("err2 ============")
				fmt.Println(err)
				//return
			}
			fmt.Println(res)
			logger.Info("队列创建消费者", res)
		}
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
}

// DelStreamGroup todo
func (q *Queue) DelStreamGroup(name string) {
	for _, stream := range q.streamList {
		if stream.Name != name {
			continue
		}
		for _, group := range stream.GroupList {
			q.Client.XGroupDestroy(context.Background(), stream.Name, group.Name).Result()
		}
	}
}

func (q *Queue) Info(queueName string) {
	//queueName = q.GlobalName + queueName
	//result, err := q.Client.XInfoGroups(context.Background(), queueName).Result()
	//logger.Info(queueName)
	//fmt.Println(result)
	//fmt.Println(err)
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

func (q *Queue) Loop() {

}
