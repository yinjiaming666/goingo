package queue

var HookMap = map[HookeName]*HookFunc{
	CallbackSuccess: &callbackSuccessFunc,
	PopSuccess:      &popSuccessFunc,
}

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

var AckMsgFail HookeName = "ack_msg_fail"

type HookFunc func(stream Stream, data map[string]any) *HookResult

var callbackSuccessFunc HookFunc = func(stream Stream, data map[string]any) *HookResult {
	//fmt.Println(hook.GetValue("Msg"))
	//logger.System("QUEUE CALLBACK SUCCESS", "Msg", hook.GetValue("Msg"))
	return &HookResult{}
}

var popSuccessFunc HookFunc = func(stream Stream, data map[string]any) *HookResult {
	return &HookResult{}
}

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
