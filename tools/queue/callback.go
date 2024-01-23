package queue

type CallbackFunc func(msg *Msg) *CallbackResult

var CallbackMap = map[string]*CallbackFunc{}

func RegisterCallback(name string, execFunc *CallbackFunc) {
	CallbackMap[name] = execFunc
}

// CallbackResult todo
type CallbackResult struct {
	Err      error
	Msg      string
	Code     int
	BackData interface{}
}
