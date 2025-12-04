package middleware

import (
	logic2 "app/internal/logic"
	"app/internal/model"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckJwt() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "请上传jwt", Body: nil}).Response()
		}
		data, err := logic2.TokenLogicInstance.CheckJwt(token)
		if err != nil {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "jwt解析失败", Body: map[string]any{"err": err.Error()}}).Response()
		}

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

			rolesGroup := new(model.RolesGroup)
			rolesGroup.Id = user.RolesGroupId
			rolesGroup.GetRolesGroup()
			auth := logic2.NewAdminAuth(user.Id, user.Pid, rolesGroup, isSuper)
			auth.Name = user.Name
			auth.Avatar = user.Avatar
			auth.Cache()
			c.Set(string(jwt.AdminJwtType), auth.Id) // c.Set() c.Get 跨中间件取值
			c.Next()
			break
		case jwt.IndexJwtType:
			user := logic2.UserLogicInstance.LoadUser(data.Uid)
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
		adminId := c.GetUint(string(jwt.AdminJwtType))
		admin := logic2.GetAdminAuth(adminId)
		fmt.Println(admin)
		fmt.Println(c.Request.URL.Path)
		menu := model.Roles{}
		menu.SearchByPath(c.Request.URL.Path)
		if menu.Id == 0 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "BackendAuth 未查询到权限", Body: nil}).Response()
		}
		checkId, _ := conv.Conv[uint](menu.Id)
		has, _ := conv.InSlice[uint](admin.RolesIds, checkId)
		if has == -1 {
			(&resp.JsonResp{Code: resp.ReAuthFail, Message: "BackendAuth 无权限访问", Body: nil}).Response()
		}
		c.Next()
	}
}
