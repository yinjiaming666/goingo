package middleware

import (
	"app/tools/logger"
	"app/tools/resp"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CatchErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			e := recover()
			if e != nil {
				switch e.(type) {
				case error:
					// 捕获错误异常
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
