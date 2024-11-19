package queue

var HookMap = map[HookeName]*HookFunc{}

type HookeName string

// PushSuccess 队列放入数据事件
var PushSuccess HookeName = "push_success"

// PopSuccess 队列取出数据事件
var PopSuccess HookeName = "pop_success"

// CallbackSuccess 执行回调成功事件
var CallbackSuccess HookeName = "callback_success"

// CallbackFail 执行回调失败事件
var CallbackFail HookeName = "callback_fail"

// UndefinedCallback 未定义的 callback 事件
var UndefinedCallback HookeName = "undefined_callback"

// AckMsgFail ack 消息失败事件
var AckMsgFail HookeName = "ack_msg_fail"

// WorkStartFail worker 启动失败
var WorkStartFail HookeName = "work_start_fail"

type HookFunc func(stream Stream, data map[string]any) *HookResult

type Hook struct {
	name *HookeName
	data map[string]any
}

type HookResult struct {
	Err      error
	Msg      string
	Code     int
	BackData interface{}
}

func RegisterHook(name HookeName, handel *HookFunc) {
	HookMap[name] = handel
}
