package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goingo/tools/logger"
	"strconv"
	"time"
)

var Client *redis.Client
var GlobalName string
var streamList = make(map[string]Stream)

type SType string

// Normal 普通队列
var Normal SType = "normal"

// Delay 延时队列
var Delay SType = "delay"

type Stream interface {
	Loop()
	Name() string
	SetName(string)
	SetFullName(string)
	FullName() string
	HandelGroup() *XGroup
	SetHandelGroup(*XGroup)
	Hook() chan *Hook
	SetHook(chan *Hook)
}

// NormalStream 普通队列
type NormalStream struct {
	name        string
	fullName    string  // redis 里存的名字
	handelGroup *XGroup // 用来执行的消费者组
	hook        chan *Hook
}

func (n *NormalStream) Loop() {
	for _, stream := range streamList {
		for _, hc := range stream.HandelGroup().ConsumerList {
			go hc.work(stream.Hook())
		}
	}

	// 用于执行钩子
	go n.listenHook()
}

func (n *NormalStream) listenHook() {
	for {
		select {
		case hook := <-n.hook:
			fun, ok := HookMap[*hook.name]
			if ok {
				_ = (*fun)(n, hook.data)
			}
			break
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (n *NormalStream) Name() string {
	return n.name
}

func (n *NormalStream) SetName(name string) {
	n.name = name
}

func (n *NormalStream) FullName() string {
	return n.fullName
}

func (n *NormalStream) SetFullName(name string) {
	n.fullName = name
}

func (n *NormalStream) HandelGroup() *XGroup {
	return n.handelGroup
}

func (n *NormalStream) SetHandelGroup(group *XGroup) {
	n.handelGroup = group
}

func (n *NormalStream) Hook() chan *Hook {
	return n.hook
}

func (n *NormalStream) SetHook(h chan *Hook) {
	n.hook = h
}

// DelayStream 延时队列
type DelayStream struct {
	name        string
	fullName    string  // redis 里存的名字
	handelGroup *XGroup // 用来执行的消费者组
	hook        chan *Hook
}

func (d *DelayStream) Loop() {
}

func (d *DelayStream) Name() string {
	return d.name
}

func (d *DelayStream) SetName(name string) {
	d.name = name
}

func (d *DelayStream) FullName() string {
	return d.fullName
}

func (d *DelayStream) SetFullName(name string) {
	d.fullName = name
}

func (d *DelayStream) HandelGroup() *XGroup {
	return d.handelGroup
}

func (d *DelayStream) SetHandelGroup(group *XGroup) {
	d.handelGroup = group
}

func (d *DelayStream) Hook() chan *Hook {
	return d.hook
}

func (d *DelayStream) SetHook(hook chan *Hook) {
	d.hook = hook
}

// XGroup 消费组
type XGroup struct {
	streamName   string
	name         string
	start        string
	ConsumerList []Consumer
}

func (g *XGroup) GetPending() (*redis.XPending, error) {
	stream, ok := streamList[g.streamName]
	if !ok {
		return nil, errors.New("not found stream")
	}
	result, err := Client.XPending(context.Background(), stream.Name(), g.name).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

type Consumer interface {
	work(chan *Hook)
	Name() string
	SetName(string)
	GroupName() string
	SetGroupName(string)
	StreamName() string
	SetStreamName(string)
	Callback() CallbackFunc
	SetCallback(CallbackFunc)
}

// HandelConsumer 消费者
type HandelConsumer struct {
	name       string
	groupName  string
	streamName string
	callback   CallbackFunc
}

func (c *HandelConsumer) Callback() CallbackFunc {
	return c.callback
}

func (c *HandelConsumer) SetCallback(exec CallbackFunc) {
	c.callback = exec
}

func (c *HandelConsumer) Name() string {
	return c.name
}

func (c *HandelConsumer) SetName(name string) {
	c.name = name
}

func (c *HandelConsumer) GroupName() string {
	return c.groupName
}

func (c *HandelConsumer) SetGroupName(groupName string) {
	c.groupName = groupName
}

func (c *HandelConsumer) StreamName() string {
	return c.streamName
}

func (c *HandelConsumer) SetStreamName(streamName string) {
	c.streamName = streamName
}

func (c *HandelConsumer) work(hook chan *Hook) {
	for {
		result, err := Client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    c.groupName,
			Consumer: c.name,
			Streams:  []string{c.streamName, ">"},
			Count:    1,
			Block:    0,
		}).Result()
		if err != nil {
			// todo hook
		}
		for _, xStream := range result {
			ml := ParseMsg(xStream.Messages)
			for _, msg := range ml {
				hook <- &Hook{
					name: &PopSuccess,
					data: map[string]any{
						"consumer": c.name,
						"msg":      msg,
					},
				}
				callbackResult := (c.Callback())(msg)
				if callbackResult.Err != nil {
					hook <- &Hook{
						name: &CallbackFail,
						data: map[string]any{
							"consumer":     c.name,
							"callback_res": callbackResult,
							"msg":          msg,
						},
					}
					continue
				}
				ack, err := Client.XAck(context.Background(), c.streamName, c.groupName, msg.Id).Result()
				if err != nil {
					hook <- &Hook{
						name: &CallbackSuccess,
						data: map[string]any{
							"consumer":     c.name,
							"callback_res": AckMsgFail,
							"msg":          msg,
							"ack":          ack,
						},
					}
					continue
				} else {
					hook <- &Hook{
						name: &CallbackSuccess,
						data: map[string]any{
							"consumer":     c.name,
							"callback_res": callbackResult,
							"msg":          msg,
						},
					}
					continue
				}
			}
		}
		//time.Sleep(1 * time.Second)
	}
}

func Init(name string, client *redis.Client) {
	Client = client
	GlobalName = name
}

func generateFullStreamName(name string, sType SType) string {
	return GlobalName + ":" + string(sType) + ":" + name
}

func CreateStream(stream Stream) error {
	if stream.Name() == "" {
		return errors.New("empty stream name")
	}

	if _, ok := streamList[stream.Name()]; ok {
		return errors.New("repeat stream")
	}

	stream.SetFullName(generateFullStreamName(stream.Name(), StreamType(stream)))
	stream.SetHook(make(chan *Hook))

	if stream.HandelGroup() == nil {
		stream.SetHandelGroup(&XGroup{
			streamName: stream.FullName(),
			name:       "name",
			start:      "$", // 指定从最后一条开始读取
			ConsumerList: []Consumer{
				&HandelConsumer{
					name: "handel1",
				},
				&HandelConsumer{
					name: "handel2",
				},
			},
		})
	}
	// 创建消费组时如果指定的 stream 不存在会报错。增加参数 MKSTREAM ，可以在 stream 不存在时自动创建它
	res, err := Client.XGroupCreateMkStream(context.Background(), stream.FullName(), stream.HandelGroup().name, stream.HandelGroup().start).Result()
	if err != nil {
		// todo
	}
	logger.Info("队列：" + stream.FullName() + "创建执行消费者组" + res)
	for k, consumer := range stream.HandelGroup().ConsumerList {
		if consumer.Name() == "" {
			return errors.New("empty consumer name")
		}
		stream.HandelGroup().ConsumerList[k].SetStreamName(stream.FullName())
		stream.HandelGroup().ConsumerList[k].SetGroupName(stream.HandelGroup().name)
		_, err = Client.XGroupCreateConsumer(context.Background(), stream.FullName(), stream.HandelGroup().name, consumer.Name()).Result()
		if err != nil {
			// todo
		}
		consumer.SetCallback(func(msg *Msg) *CallbackResult {
			fun, ok := CallbackMap[msg.CallbackName]
			if !ok {
				stream.Hook() <- &Hook{
					name: &UndefinedCallback,
					data: map[string]any{
						"msg": msg,
					},
				}
				return &CallbackResult{
					Err:      errors.New("undefined callback"),
					Msg:      "undefined callback",
					Code:     1,
					BackData: nil,
				}
			} else {
				return (*fun)(msg)
			}
		})
		logger.Info("队列：" + stream.FullName() + "执行消费者组创建消费者：" + consumer.Name())
	}

	streamList[stream.Name()] = stream
	EchoInfo(stream.Name())
	return nil
}

// EchoInfo 输出 stream 全部信息
func EchoInfo(streamName string) {
	s, ok := streamList[streamName]
	if !ok {
		fmt.Println("未查询到队列")
		return
	}

	stream, err := Client.XInfoStream(context.Background(), s.FullName()).Result()
	if err != nil {
		fmt.Printf("队列 %s 获取失败：%v \n", s.FullName(), err)
		return
	}
	fmt.Printf(">队列名称：%s \n", s.FullName())
	fmt.Printf(">队列类型：%s \n", string(StreamType(s)))
	fmt.Printf(">队列长度：%d \n", stream.Length)
	fmt.Printf(">FirstEntry：%v \n", stream.FirstEntry)
	fmt.Printf(">LastEntry：%v \n", stream.LastEntry)
	fmt.Printf(">RecordedFirstEntryID：%s \n", stream.RecordedFirstEntryID)

	groups, err := Client.XInfoGroups(context.Background(), s.FullName()).Result()
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
		consumers, _ := Client.XInfoConsumers(context.Background(), s.FullName(), group.Name).Result()
		for _, consumer := range consumers {
			fmt.Printf(">        消费者名称：%s \n", consumer.Name)
			fmt.Printf(">        Pending：%d \n", consumer.Pending)
			fmt.Printf(">        Idle：%s \n", consumer.Idle)
			fmt.Printf(">        Inactive：%s \n", consumer.Inactive)
		}
	}
	return
}

func StreamType(stream Stream) SType {
	switch stream.(type) {
	case *NormalStream:
		return Normal
	case *DelayStream:
		return Delay
	default:
		return ""
	}
}

func Push(queueName, c, callback string, data map[string]interface{}) (string, error) {
	stream := streamList[queueName]
	var msg = Msg{
		C:            c,
		CallbackName: callback,
		Data:         data,
	}
	var b = &redis.XAddArgs{
		Stream: stream.FullName(),
		MaxLen: 0,
		ID:     "",
		Values: map[string]interface{}{
			"data": msg,
		},
	}
	result, err := Client.XAdd(context.Background(), b).Result()
	if err != nil {
		return "", err
	}
	msg.Id = result
	stream.Hook() <- &Hook{
		name: &PushSuccess,
		data: map[string]any{
			"msg": &msg,
		},
	}
	return result, err
}

func PushDaly(body map[string]interface{}, queueName string, sType SType, second int) (string, error) {
	id := time.Now().UnixMilli()
	id += int64(second) * int64(time.Millisecond)
	fmt.Println(strconv.Itoa(int(id)))

	var b = &redis.XAddArgs{
		Stream: generateFullStreamName(queueName, sType),
		MaxLen: 0,
		ID:     strconv.Itoa(int(id)),
		Values: body,
	}
	return Client.XAdd(context.Background(), b).Result()
}
