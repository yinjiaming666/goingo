package middleware

import (
	"app/internal/global"
	"app/tools/conv"
	"app/tools/resp"
	"github.com/gin-gonic/gin"
)

func FilterIp(allowIp []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if k, _ := conv.InSlice(allowIp, "*"); k >= 0 {
			// 允许所有 ip 访问
			c.Next()
			return
		}

		ip := c.ClientIP()
		if ip == "::1" || ip == "localhost" {
			ip = global.LocalIp
		}
		if k, _ := conv.InSlice(allowIp, ip); k < 0 {
			(&resp.JsonResp{Code: resp.ReIllegalIp, Message: "illegal IP", Body: map[string]any{"ip": ip}}).Response()
		}
		c.Next()
	}
}
