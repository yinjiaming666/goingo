# Goingo

基于 Gin + Gorm 整合的开发框架，用于快速构建 API 服务

## 使用技术

- 路由，中间件 [Gin](https://gin-gonic.com/zh-cn/)
- model [Gorm](https://gorm.io/zh_CN/)
- 配置文件解析 [viper](https://github.com/spf13/viper/)
- [jwt](https://github.com/golang-jwt/jwt/)
- [redis](https://redis.uptrace.dev/zh/)

## 目录结构

```
├── config              // 项目配置文件
│   ├── dev.ini         // 开发环境
│   ├── prod.ini        // 测试环境
│   └── server.ini      // 服务器信息
├── deploy              // 打包上传到正式环境
│   ├── deploy.go
│   ├── deploy.sh
│   └── run.sh
├── internal            // 业务代码
│   ├── logic           // 业务逻辑
│   ├── middleware      // 中间件
│   ├── model           // 模型
│   ├── router          // 路由
│   └── server          // 接口
├── log                 // 运行日志
├── main.go             // 入口文件
└── tools               // 通用工具
    ├── conv
    ├── jwt            
    ├── key_utils
    ├── logger          // 日志
    ├── queue           // 队列
    ├── random
    ├── resp            // 响应
    └── utils.go
```

## 运行

```shell
go run main.go -mode=dev
# 运行参数
# -mode=dev    运行测试环境 dev.ini
# -mode=prod   运行正式环境 prod.ini
# -initDb=true 根据结构体初始化数据库
```

## 队列

> 队列是基于 redis 的 stream 实现的 <br> 延时队列是基于 redis 的 zset 实现的

### 初始化队列

```
queue.Init("goingo-queue", model.RedisClient)

// 延时队列
stream := &queue.DelayStream{}
stream.SetName("default")
err := stream.Create()  //（redis key name goingo-queue:delay:default）
if err != nil {
    fmt.Println(err.Error())
    return
}
go stream.Loop()

// 消息队列
stream = &queue.NormalStream{}
stream.SetName("default")
err := stream.Create()  //（redis key name goingo-queue:normal:default）
if err != nil {
    fmt.Println(err.Error())
    return
}
go stream.Loop()
```

### 队列投入数据

#### 消息队列

> `queue.Push(队列名称, 回调名称, map[string]interface{}{"name": "张三", "age": 19})`

#### 延时队列

> `queue.PushDelay(队列名称, 回调名称, map[string]interface{}{"name": "张三", "age": 19}, 延时秒数)`

### 回调与钩子

#### 注册回调

```
var pF queue.CallbackFunc = func(msg *queue.Msg) *queue.CallbackResult {
	// 业务逻辑
	return &queue.CallbackResult{
		Err:      nil,
		Msg:      "",
		Code:     0, // 0 成功，1 失败
		BackData: nil,
	}
}
queue.RegisterCallback("test", &pF)
```

#### 注册钩子

```
var u queue.HookFunc = func(stream queue.Stream, data map[string]any) *queue.HookResult {
	_, ok := data["msg"]
	if !ok {
		return &queue.HookResult{
			Err:      errors.New("nil msg"),
			Msg:      "nil msg",
			Code:     1,
			BackData: nil,}
	}
	msg := data["msg"].(*queue.Msg)
	_, ok = data["consumer"]
	if !ok {
		return &queue.HookResult{
			Err:      errors.New("nil consumer"),
			Msg:      "nil consumer",
			Code:     1,
			BackData: nil,}
	}
	consumer := data["consumer"].(string)
	logger.System("CALLBACK MSG", "Msg", msg.Id, "consumer", consumer)
	queue.Client.XDel(context.Background(), stream.FullName(), msg.Id)
	return &queue.HookResult{
		Err:      nil,
		Msg:      "success",
		Code:     0,
		BackData: nil,
	}
}
queue.RegisterHook(queue.CallbackSuccess, &u)
```

#### 钩子事件列表

<ul>
    <li>PushSuccess 队列放入数据事件</li>
    <li>PopSuccess 队列取出数据事件</li>
    <li>CallbackSuccess 执行回调成功事件</li>
    <li>CallbackFail 执行回调失败事件</li>
    <li>UndefinedCallback 未定义的 callback 事件</li>
</ul>

### 队列完整示例

```go
package main

import (
	"fmt"
	"goingo/internal/model"
	"goingo/tools/queue"
)

func main() {
	logger.InitLog()
	model.InitRedis(&model.RedisConf{
		Ip:         "192.168.110.177",
		Port:       "63792",
		GlobalName: "goingo-queue",
	})
	queue.Init("goingo-queue", model.RedisClient)
	stream := &queue.NormalStream{}
	stream.SetName("default")
	err := stream.Create() // 初始化创建队列（redis key name goingo-queue:normal:default）
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 队列投入数据，callbackName 需要通过 RegisterHook 注册回调
	queue.Push("default", "test", map[string]interface{}{"name": "张三", "age": 19})

	// 注册回调
	var pF queue.CallbackFunc = func(msg *queue.Msg) *queue.CallbackResult {
		// 业务逻辑
		return &queue.CallbackResult{
			Err:      nil,
			Msg:      "",
			Code:     0, // 0 成功，1 失败
			BackData: nil,
		}
	}
	queue.RegisterCallback("test", &pF)

	// 注册钩子
	var u queue.HookFunc = func(stream queue.Stream, data map[string]any) *queue.HookResult {
		_, ok := data["msg"]
		if !ok {
			return &queue.HookResult{
				Err:      errors.New("nil msg"),
				Msg:      "nil msg",
				Code:     1,
				BackData: nil,
			}
		}
		msg := data["msg"].(*queue.Msg)

		_, ok = data["consumer"]
		if !ok {
			return &queue.HookResult{
				Err:      errors.New("nil consumer"),
				Msg:      "nil consumer",
				Code:     1,
				BackData: nil,
			}
		}
		consumer := data["consumer"].(string)
		logger.System("CALLBACK MSG", "Msg", msg.Id, "consumer", consumer)
		return &queue.HookResult{
			Err:      nil,
			Msg:      "success",
			Code:     0,
			BackData: nil,
		}
	}
	queue.RegisterHook(queue.CallbackSuccess, &u)
	stream.Loop()
}
```

## 打包上传到服务器

```shell
cd deploy && go run deploy.go
```
