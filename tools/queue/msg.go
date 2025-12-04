package queue

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Msg struct {
	C            string // 保留字段
	CallbackName string
	Id           string
	MsgType      SType
	Data         string // json 字符串
}

func (s Msg) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func ParseMsg(list []redis.XMessage) []*Msg {
	var l []*Msg
	for _, XMessage := range list {
		if v, ok := XMessage.Values["data"]; ok {
			m := Json2Msg(v.(string))
			m.Id = XMessage.ID
			l = append(l, m)
		}
	}
	return l
}

func Json2Msg(jsonStr string) *Msg {
	var m Msg
	_ = json.Unmarshal([]byte(jsonStr), &m)
	return &m
}
