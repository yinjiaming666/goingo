package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	global "goingo/internal"
	"goingo/tools/conv"
	"goingo/tools/resp"
)

func FilterIp(allowIp []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if k, _ := conv.InSlice(allowIp, "*"); k >= 0 {
			// 允许所有 ip 访问
			c.Next()
			return
		}

		ip := c.ClientIP()
		fmt.Println(ip)
		if ip == "::1" || ip == "localhost" {
			ip = global.LocalIp
		}
		if k, _ := conv.InSlice(allowIp, ip); k < 0 {
			// 允许所有 ip 访问
			resp.Resp(resp.ReIllegalIp, "illegal IP", map[string]any{"ip": ip})
			return
		}
		c.Next()
	}
}
