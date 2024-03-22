package resp

import "github.com/gin-gonic/gin"

type Code uint

const ReSuccess Code = 200
const ReFail Code = 400
const ReAuthFail Code = 401
const ReIllegalIp Code = 402
const ReError Code = 500

type Response interface {
	GetCode() Code
	GetBody() any
	GetMsg() string
	GetHttpCode() int
	SetHttpCode(code int)
	Response()
}

type JsonResp struct {
	Code     Code
	Message  string
	Body     any
	httpCode int
}

func (j *JsonResp) GetCode() Code {
	return j.Code
}

func (j *JsonResp) GetBody() any {
	return j.Body
}

func (j *JsonResp) GetMsg() string {
	return j.Message
}

func (j *JsonResp) GetHttpCode() int {
	if j.httpCode == 0 {
		return 200
	}
	return j.httpCode
}

func (j *JsonResp) SetHttpCode(code int) {
	j.httpCode = code
}

func (j *JsonResp) Response() {
	panic(j)
}

type StringResp struct {
	Code     Code
	Message  string
	Body     string
	httpCode int
}

func (s *StringResp) GetCode() Code {
	return s.Code
}

func (s *StringResp) GetBody() any {
	return s.Body
}

func (s *StringResp) GetMsg() string {
	return s.Message
}

func (s *StringResp) Response() {
	panic(s)
}

func (s *StringResp) GetHttpCode() int {
	if s.httpCode == 0 {
		return 200
	}
	return s.httpCode
}

func (s *StringResp) SetHttpCode(code int) {
	s.httpCode = code
}

type XmlResp struct {
	Code     Code
	Message  string
	Body     any // struct
	httpCode int
}

func (x *XmlResp) GetCode() Code {
	return x.Code
}

func (x *XmlResp) GetBody() any {
	return x.Body
}

func (x *XmlResp) GetMsg() string {
	return x.Message
}

func (x *XmlResp) Response() {
	panic(x)
}

func (x *XmlResp) GetHttpCode() int {
	if x.httpCode == 0 {
		return 200
	}
	return x.httpCode
}

func (x *XmlResp) SetHttpCode(code int) {
	x.httpCode = code
}

func HandelResponse(r Response, c *gin.Context) {
	switch r.(type) {
	case *StringResp:
		c.Abort()
		c.String(r.GetHttpCode(), "%s", r.GetBody())
		c.Next()
		break
	case *JsonResp:
		c.Abort()
		c.JSON(r.GetHttpCode(), gin.H{
			"Code": r.GetCode(),
			"msg":  r.GetMsg(),
			"data": r.GetBody(),
		})
		c.Next()
		break
	case *XmlResp:
		c.Abort()
		c.JSON(r.GetHttpCode(), r.GetBody())
		c.Next()
		break
	}
}
