package resp

type Code uint

const ReSuccess Code = 200
const ReFail Code = 400
const ReError Code = 500

type Response struct {
	Code    Code
	Message string
	Data    interface{}
}

func Resp(code Code, msg string, obj interface{}) {
	panic(&Response{
		Code:    code,
		Message: msg,
		Data:    obj,
	})
}
