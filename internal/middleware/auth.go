package middleware

import (
	"github.com/gin-gonic/gin"
	logic2 "goingo/internal/logic"
	"goingo/internal/model"
	"goingo/tools/jwt"
	"goingo/tools/resp"
)

func CheckJwt() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			resp.Resp(resp.ReFail, "请上传jwt", nil)
		}
		// golang 变量作用域跟 js 的 let 类似，for if switch 中声明的变量不能拿到外面去用
		tl := logic2.TokenLogic{}
		data := tl.CheckJwt(token)

		switch data.Type {
		case jwt.AdminJwtType:
			user := &model.Admin{
				Id: data.Uid,
			}
			user = user.GetAdmin()
			if user.Id <= 0 {
				resp.Resp(resp.ReFail, "未查询到用户", nil)
			}
			c.Set(string(jwt.AdminJwtType), user) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		case jwt.IndexJwtType:
			userLogic := logic2.UserLogic{}
			user := userLogic.LoadUser(data.Uid)
			if user.Id <= 0 {
				resp.Resp(resp.ReFail, "未查询到用户", nil)
			}
			c.Set(string(jwt.IndexJwtType), user) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		}

	}
}
