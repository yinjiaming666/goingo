package middleware

import (
	logic2 "app/internal/logic"
	"app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func CheckJwt() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "请上传jwt", Body: nil}).Response()
		}
		tl := logic2.TokenLogic{}
		data := tl.CheckJwt(token)

		switch data.Type {
		case jwt.AdminJwtType:
			user := &model.Admin{
				Id: data.Uid,
			}
			user = user.GetAdmin()
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Message: "未查询到用户", Body: nil}).Response()
			}
			isSuper := user.IsSuper == 1

			groupIds := (new(model.AdminRolesGroup)).GetGroupIdsByAdminId(user.Id)

			auth := logic2.NewAdminAuth(user.Id, user.Pid, groupIds, isSuper)
			auth.Name = user.Name
			auth.Avatar = user.Avatar
			auth.Cache()
			c.Set(string(jwt.AdminJwtType), auth.Id) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		case jwt.IndexJwtType:
			userLogic := logic2.UserLogic{}
			user := userLogic.LoadUser(data.Uid)
			if user.Id <= 0 {
				(&resp.JsonResp{Code: resp.ReAuthFail, Message: "未查询到用户", Body: nil}).Response()
			}
			c.Set(string(jwt.IndexJwtType), user)
			c.Next()
			break
		}

	}
}

// BackendAuth 管理后台鉴权
func BackendAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		var admin model.Admin
		a, _ := c.Get(string(jwt.AdminJwtType))
		if t, ok := a.(model.Admin); !ok {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "BackendAuth 解析错误", Body: nil}).Response()
		} else {
			admin = t
		}
		fmt.Println(admin)
		fmt.Println(c.Request.URL.Path)
		menu := model.Roles{}
		menu.SearchByPath(c.Request.URL.Path)
		if menu.Id == 0 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "BackendAuth 未查询到权限", Body: nil}).Response()
		}
		roleList := strings.Split(admin.RolesGroupIds, ",")
		checkId, _ := conv.Conv[string](menu.Id)
		_, ok := conv.InSlice[string](roleList, checkId)
		if ok == "" {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "BackendAuth 无权限访问", Body: nil}).Response()
		}
		c.Next()
	}
}
