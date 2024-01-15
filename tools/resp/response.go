package resp

const ReSuccess = 200
const ReFail = 400

type Response struct {
	Code    int
	Message string
	Data    interface{}
}

func Resp(code int, msg string, obj interface{}) {
	panic(&Response{
		Code:    code,
		Message: msg,
		Data:    obj,
	})
}
