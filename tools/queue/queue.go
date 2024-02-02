package queue

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"goingo/tools/logger"
	"goingo/tools/random"
	"strconv"
	"time"
)

var Client *redis.Client
var GlobalName string
var streamList = make(map[string]Stream)

type SType string

// Normal 消息队列
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
	Create() error
}

// NormalStream 消息队列
type NormalStream struct {
	name        string
	fullName    string  // redis 里存的名字
	handelGroup *XGroup // 用来执行的消费者组
	hook        chan *Hook
}

func (n *NormalStream) Create() error {
	if n.Name() == "" {
		return errors.New("empty stream name")
	}

	if _, ok := streamList[string(Normal)+"-"+n.Name()]; ok {
		return errors.New("repeat stream:'" + string(Normal) + "-" + n.Name() + "'")
	}

	n.SetFullName(generateFullStreamName(n.Name(), Normal))
	n.SetHook(make(chan *Hook))

	if n.HandelGroup() == nil {
		n.SetHandelGroup(&XGroup{
			streamName: n.FullName(),
			name:       "name",
			start:      "$", // 指定从最后一条开始读取
			ConsumerList: []Consumer{
				&NormalConsumer{
					name: "handel1",
				},
				&NormalConsumer{
					name: "handel2",
				},
			},
		})
	}
	// 创建消费组时如果指定的 stream 不存在会报错。增加参数 MKSTREAM ，可以在 stream 不存在时自动创建它
	res, err := Client.XGroupCreateMkStream(context.Background(), n.FullName(), n.HandelGroup().name, n.HandelGroup().start).Result()
	if err != nil {
		// todo
	}
	logger.Info("队列：" + n.FullName() + "创建执行消费者组" + res)
	for k, consumer := range n.HandelGroup().ConsumerList {
		if consumer.Name() == "" {
			return errors.New("empty consumer name")
		}
		n.HandelGroup().ConsumerList[k].SetStreamName(n.FullName())
		n.HandelGroup().ConsumerList[k].SetGroupName(n.HandelGroup().name)
		_, err = Client.XGroupCreateConsumer(context.Background(), n.FullName(), n.HandelGroup().name, consumer.Name()).Result()
		if err != nil {
			// todo
		}
		consumer.SetCallback(func(msg *Msg) *CallbackResult {
			fun, ok := CallbackMap[msg.CallbackName]
			if !ok {
				n.Hook() <- &Hook{
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
		logger.Info("队列：" + n.FullName() + "执行消费者组创建消费者：" + consumer.Name())
	}

	streamList[string(Normal)+"-"+n.Name()] = n
	//EchoInfo(n.Name())
	return nil
}

func (n *NormalStream) Loop() {
	for _, hc := range n.HandelGroup().ConsumerList {
		go hc.work(n.Hook())
	}
	// 用于执行钩子
	go listenHook(n)
	select {}
}

func listenHook(s Stream) {
	for {
		select {
		case hook := <-s.Hook():
			fun, ok := HookMap[*hook.name]
			if ok {
				_ = (*fun)(s, hook.data)
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

func (d *DelayStream) Create() error {
	if d.Name() == "" {
		return errors.New("empty stream name")
	}
	if _, ok := streamList[string(Delay)+"-"+d.Name()]; ok {
		return errors.New("repeat stream:'" + string(Delay) + "-" + d.Name() + "'")
	}
	d.SetFullName(generateFullStreamName(d.Name(), Delay))
	d.SetHandelGroup(&XGroup{
		streamName: d.FullName(),
		name:       "group",
		start:      "",
		ConsumerList: []Consumer{
			&DelayConsumer{
				name:       "handel1",
				streamName: d.FullName(),
				groupName:  "group",
			},
		},
	})
	d.HandelGroup().ConsumerList[0].SetCallback(
		func(msg *Msg) *CallbackResult {
			fun, ok := CallbackMap[msg.CallbackName]
			if !ok {
				d.Hook() <- &Hook{
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
		},
	)
	d.SetFullName(generateFullStreamName(d.Name(), Delay))
	d.SetHook(make(chan *Hook))
	streamList[string(Delay)+"-"+d.Name()] = d
	return nil
}

func (d *DelayStream) Loop() {
	// 用于执行钩子
	go listenHook(d)
	go d.handelGroup.ConsumerList[0].work(d.Hook())
	select {}
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

// NormalConsumer 消费者
type NormalConsumer struct {
	name       string
	groupName  string
	streamName string
	callback   CallbackFunc
}

func (c *NormalConsumer) Callback() CallbackFunc {
	return c.callback
}

func (c *NormalConsumer) SetCallback(exec CallbackFunc) {
	c.callback = exec
}

func (c *NormalConsumer) Name() string {
	return c.name
}

func (c *NormalConsumer) SetName(name string) {
	c.name = name
}

func (c *NormalConsumer) GroupName() string {
	return c.groupName
}

func (c *NormalConsumer) SetGroupName(groupName string) {
	c.groupName = groupName
}

func (c *NormalConsumer) StreamName() string {
	return c.streamName
}

func (c *NormalConsumer) SetStreamName(streamName string) {
	c.streamName = streamName
}

func (c *NormalConsumer) work(hook chan *Hook) {
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
						name: &AckMsgFail,
						data: map[string]any{
							"consumer":     c.name,
							"callback_res": callbackResult,
							"msg":          msg,
							"ack":          ack,
						},
					}
					continue
				}
				hook <- &Hook{
					name: &CallbackSuccess,
					data: map[string]any{
						"consumer":     c.name,
						"callback_res": callbackResult,
						"msg":          msg,
					},
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

type DelayConsumer struct {
	name       string
	groupName  string
	streamName string
	callback   CallbackFunc
}

func (d *DelayConsumer) Name() string {
	return d.name
}

func (d *DelayConsumer) SetName(name string) {
	d.name = name
}

func (d *DelayConsumer) GroupName() string {
	return d.groupName
}

func (d *DelayConsumer) SetGroupName(groupName string) {
	d.groupName = groupName
}

func (d *DelayConsumer) StreamName() string {
	return d.streamName
}

func (d *DelayConsumer) SetStreamName(streamName string) {
	d.streamName = streamName
}

func (d *DelayConsumer) Callback() CallbackFunc {
	return d.callback
}

func (d *DelayConsumer) SetCallback(callback CallbackFunc) {
	d.callback = callback
}

func (d *DelayConsumer) work(hook chan *Hook) {
	for {
		now := time.Now().Unix()
		var nowStr string
		nowStr = strconv.FormatInt(now, 10)
		result, err := Client.ZRangeByScore(context.Background(), d.StreamName(), &redis.ZRangeBy{
			Min: "0",
			Max: nowStr,
		}).Result()
		if err != nil {
			// todo
		}
		for _, member := range result {
			msg := Json2Msg(member)
			hook <- &Hook{
				name: &PopSuccess,
				data: map[string]any{
					"consumer": d.name,
					"msg":      msg,
				},
			}
			callbackResult := (d.Callback())(msg)
			if callbackResult.Err != nil {
				hook <- &Hook{
					name: &CallbackFail,
					data: map[string]any{
						"consumer":     d.name,
						"callback_res": callbackResult,
						"msg":          msg,
					},
				}
				continue
			}
			// todo

			i, err := Client.ZRem(context.Background(), d.StreamName(), member).Result()
			if err != nil {
				hook <- &Hook{
					name: &AckMsgFail,
					data: map[string]any{
						"consumer":     d.name,
						"callback_res": callbackResult,
						"msg":          msg,
						"ack":          i,
					},
				}
				continue
			}
			hook <- &Hook{
				name: &CallbackSuccess,
				data: map[string]any{
					"consumer":     d.name,
					"callback_res": callbackResult,
					"msg":          msg,
				},
			}
			continue
		}
		time.Sleep(1 * time.Second)
	}
}

func Init(name string, client *redis.Client) {
	Client = client
	GlobalName = name
}

func generateFullStreamName(name string, sType SType) string {
	return GlobalName + ":" + string(sType) + ":" + name
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

func Push(queueName, callback string, data map[string]interface{}) (string, error) {
	stream := streamList[string(Normal)+"-"+queueName]
	var msg = Msg{
		C:            callback,
		MsgType:      Normal,
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

func PushDelay(queueName, callback string, data map[string]any, second int) (int64, error) {
	var score int64
	score = time.Now().Unix() + int64(second)

	var msg = Msg{
		C:            callback,
		MsgType:      Delay,
		CallbackName: callback,
		Data:         data,
		Id:           strconv.FormatInt(score, 10) + "-" + strconv.Itoa(random.Number(10000, 99999)),
	}
	stream := streamList[string(Delay)+"-"+queueName]
	return Client.ZAdd(context.Background(), stream.FullName(), redis.Z{
		Member: msg,
		Score:  float64(score),
	}).Result()
}
