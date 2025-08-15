package message

type C string
type M string

var CDefault C = "default"
var MHandelMoney M = "handelMoney"

type Message struct {
	C     C      `json:"c"`
	M     M      `json:"m"`
	Data  string `json:"data"` // json 字符串
	JobId uint64 `json:"job_id"`
}

type HandelMoneyMsg struct {
	T          uint8   `json:"t"`
	Num        float64 `json:"num"`
	Uid        int     `json:"uid"`
	LogType    uint8   `json:"logType"`
	AppendData string  `json:"appendData"`
	MoneyType  uint8   `json:"moneyType"` // 金币类型(保留字段)
}
