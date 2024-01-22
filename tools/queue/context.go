package queue

type Context interface {
	MsgId() string
	SetMsgId(msgId string)
	Data() any
	SetData(data any)
}

type ReadGroupContext struct {
	msgId string
	data  any
}

func (c *ReadGroupContext) Data() any {
	return c.data
}

func (c *ReadGroupContext) SetData(data any) {
	c.data = data
}

func (c *ReadGroupContext) MsgId() string {
	return c.msgId
}

func (c *ReadGroupContext) SetMsgId(msgId string) {
	c.msgId = msgId
}

type CallBackSuccessContext struct {
	msgId string
	data  any
}

func (c *CallBackSuccessContext) Data() any {
	return c.data
}

func (c *CallBackSuccessContext) SetData(data any) {
	c.data = data
}

func (c *CallBackSuccessContext) MsgId() string {
	return c.msgId
}

func (c *CallBackSuccessContext) SetMsgId(msgId string) {
	c.msgId = msgId
}
