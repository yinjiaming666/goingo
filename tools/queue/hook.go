package queue

var HookMap = map[HookFuncName]*HookFunc{
	CallbackSuccess: &callbackSuccessFunc,
	PopSuccess:      &popSuccessFunc,
}

type HookFuncName int

// PushSuccess 队列放入数据事件
var PushSuccess HookFuncName = 0

// PopSuccess 队列取出数据事件
var PopSuccess HookFuncName = 1

// CallbackSuccess 执行回调成功事件
var CallbackSuccess HookFuncName = 2

// CallbackFail 执行回调失败事件
var CallbackFail HookFuncName = 3

// UndefinedCallback 未定义的 callback 事件
var UndefinedCallback HookFuncName = 4

var AckMsgFail HookFuncName = 5

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
	name *HookFuncName
	data map[string]any
}

type HookResult struct {
	Err      error
	Msg      string
	Code     int
	BackData interface{}
}

func RegisterHook(name HookFuncName, handel *HookFunc) {
	HookMap[name] = handel
}
