package queue

import "fmt"

type ExecFunc func(msg *Msg) ExecResult

var ExecFuncMap = map[string]*ExecFunc{
	"print":      &pF,
	"backup_log": &backupLog,
}

var pF ExecFunc = func(msg *Msg) ExecResult {
	fmt.Println("exec ==============")
	fmt.Println(msg)
	return ExecResult{
		err:      nil,
		msg:      "",
		code:     0,
		backData: nil,
	}
}

var backupLog ExecFunc = func(msg *Msg) ExecResult {
	fmt.Println("backupLog ==============")
	fmt.Println(msg)
	return ExecResult{
		err:      nil,
		msg:      "",
		code:     0,
		backData: nil,
	}
}

func RegisterExec(name string, execFunc *ExecFunc) {
	ExecFuncMap[name] = execFunc
}

// ExecResult todo
type ExecResult struct {
	err      error
	msg      string
	code     int
	backData interface{}
}
