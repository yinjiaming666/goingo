package resp

import (
	"app/tools/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

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
	Response(c *gin.Context)
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

func (j *JsonResp) Response(c *gin.Context) {
	addRespLog(j, c)
	c.AbortWithStatusJSON(j.GetHttpCode(), gin.H{
		"code": j.GetCode(),
		"msg":  j.GetMsg(),
		"data": j.GetBody(),
	})
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

func (s *StringResp) Response(ctx *gin.Context) {
	addRespLog(s, ctx)
	ctx.String(s.GetHttpCode(), "%s", s.GetBody())
	ctx.Abort()
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

func (x *XmlResp) Response(c *gin.Context) {
	addRespLog(x, c)
	c.XML(x.GetHttpCode(), x.GetBody())
	c.Abort()
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

func addRespLog(response Response, c *gin.Context) {
	r, _ := json.Marshal(response.GetBody())
	logger.System("Response", "method", c.Request.Method, "url", c.Request.URL.String(), "post", c.Request.PostForm, "res", map[string]any{
		"code": response.GetCode(),
		"msg":  response.GetMsg(),
		"data": string(r),
	})
}
