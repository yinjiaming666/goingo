package index_api

import (
	logic2 "app/internal/logic"
	model2 "app/internal/model"
	"app/tools"
	"app/tools/jwt"
	"app/tools/resp"
	"github.com/gin-gonic/gin"
)

// RegisterUser 注册用户
func RegisterUser(content *gin.Context) {
	nickname := content.PostForm("nickname")
	password := content.PostForm("password")

	userLogic := logic2.UserLogic{}
	search := make(map[string]any)
	search["nickname"] = nickname
	rep := userLogic.SearchUser(search)
	if rep.Id > 0 {
		(&resp.JsonResp{Code: resp.ReSuccess, Message: "当前用户已注册", Body: nil}).Response()
	}

	user := model2.User{
		Nickname: nickname,
		Password: password,
	}
	uid := user.DoRegister()
	if uid == 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "注册失败", Body: nil}).Response()
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "注册成功", Body: map[string]any{"user": user}}).Response()
}

// Login 用户登陆
func Login(content *gin.Context) {
	nickname := content.PostForm("nickname")
	password := content.PostForm("password")

	s := map[string]any{
		"nickname": nickname,
		"password": tools.Md5(password, model2.UserPwdSalt),
	}

	userInfo := logic2.UserLogic{}.SearchUser(s)
	if userInfo.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "账号或密码错误", Body: nil}).Response()
		content.Abort()
		return
	}
	data := make(map[string]interface{})

	tokenLogic := logic2.TokenLogic{}
	j, userJwt := tokenLogic.GenerateJwt(userInfo.Id, jwt.IndexJwtType, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = userInfo
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: data}).Response()
}

func LoadUser(content *gin.Context) {
	user, ok := content.Get(string(jwt.IndexJwtType))
	if !ok {
		(&resp.JsonResp{Code: resp.ReFail, Message: "未查询到用户", Body: nil}).Response()
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "success", Body: map[string]any{"user": user}}).Response()
}
