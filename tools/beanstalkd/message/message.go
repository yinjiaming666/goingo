package message

type C string
type M string

var CDefault C = "default"
var MSaveUser M = "updateUser"

type Message struct {
	C     C      `json:"c"`
	M     M      `json:"m"`
	Data  string `json:"data"` // json 字符串
	JobId uint64 `json:"job_id"`
}
