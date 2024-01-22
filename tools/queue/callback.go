package queue

import (
	"fmt"
	"goingo/tools/logger"
)

type CallbackFunc func(msg *Msg) CallbackResult

var CallbackMap = map[string]*CallbackFunc{
	"print":      &pF,
	"backup_log": &backupLog,
}

var pF CallbackFunc = func(msg *Msg) CallbackResult {
	fmt.Println("callback ==============")
	return CallbackResult{
		err:      nil,
		msg:      "",
		code:     0,
		backData: nil,
	}
}

var backupLog CallbackFunc = func(msg *Msg) CallbackResult {
	logger.System("QUEUE PUSH "+msg.Id, "data", *msg)
	return CallbackResult{
		err:      nil,
		msg:      "success",
		code:     0,
		backData: nil,
	}
}

func RegisterExec(name string, execFunc *CallbackFunc) {
	CallbackMap[name] = execFunc
}

// CallbackResult todo
type CallbackResult struct {
	err      error
	msg      string
	code     int
	backData interface{}
}
