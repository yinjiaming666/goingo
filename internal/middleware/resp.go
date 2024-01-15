package middleware

import (
	"github.com/gin-gonic/gin"
	"goingo/tools/resp"
)

func RespMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			e := recover()
			if e != nil {
				err, ok := e.(*resp.Response)
				if ok {
					c.AbortWithStatusJSON(200, gin.H{
						"code": err.Code,
						"msg":  err.Message,
						"data": err.Data,
					})
					return
				}
				c.Next()
			}
			c.Next()
		}()
		c.Next()
	}
}
