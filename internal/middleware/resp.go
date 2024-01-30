package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goingo/tools"
	"goingo/tools/logger"
	"goingo/tools/resp"
)

func RespMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			e := recover()
			if e != nil {
				switch e.(type) {
				case *resp.Response:
					// 捕获响应
					logger.Info("Response", "method", c.Request.Method, "url", c.Request.URL.String(), "post", c.Request.PostForm, "res", map[string]any{
						"code": e.(*resp.Response).Code,
						"msg":  e.(*resp.Response).Message,
						"data": e.(*resp.Response).Data,
					})
					c.AbortWithStatusJSON(200, gin.H{
						"code": e.(*resp.Response).Code,
						"msg":  e.(*resp.Response).Message,
						"data": e.(*resp.Response).Data,
					})
					c.Next()
					return
				case error:
					// 捕获错误异常
					fmt.Println(tools.PrintStackTrace(e.(error))) // 打印堆栈信息
					logger.Error(e.(error).Error(), "method", c.Request.Method, "url", c.Request.URL.String(), "post", c.Request.PostForm)
					c.AbortWithStatusJSON(200, gin.H{
						"code": resp.ReError,
						"msg":  e.(error).Error(),
						"data": map[string]any{},
					})
					c.Next()
					return
				default:
					fmt.Println("unknown recover")
					fmt.Println(e)
					c.Next()
					return
				}
			}
		}()
		c.Next()
	}
}
