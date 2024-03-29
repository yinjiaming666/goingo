package index_api

import (
	"github.com/gin-gonic/gin"
	logic2 "goingo/internal/logic"
	model2 "goingo/internal/model"
	"goingo/tools"
	"goingo/tools/jwt"
	"goingo/tools/resp"
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

	user := model2.User{
		Nickname: nickname,
		Password: tools.Md5(password, model2.UserPwdSalt),
	}

	userInfo := user.GetUserInfo()
	if userInfo.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "账号或密码错误", Body: nil}).Response()
		content.Abort()
		return
	}
	data := make(map[string]interface{})

	tokenLogic := logic2.TokenLogic{}
	j, userJwt := tokenLogic.GenerateJwt(user.Id, jwt.IndexJwtType, 0)
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
