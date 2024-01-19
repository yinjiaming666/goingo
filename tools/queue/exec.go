package queue

import "fmt"

type ExecFunc func(msg *Msg)

var ExecFuncMap = map[string]*ExecFunc{
	"print":      &pF,
	"backup_log": &backupLog,
}

var pF ExecFunc = func(msg *Msg) {
	fmt.Println("exec ==============")
	fmt.Println(msg)
}

var backupLog ExecFunc = func(msg *Msg) {
	fmt.Println("backupLog ==============")
	fmt.Println(msg)
}

func RegisterExec(name string, execFunc *ExecFunc) {
	ExecFuncMap[name] = execFunc
}
