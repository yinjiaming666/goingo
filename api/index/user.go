package index_api

import (
	"github.com/gin-gonic/gin"
	"goingo/database"
	"goingo/logic"
	"goingo/utils"
	"goingo/utils/jwt"
	"goingo/utils/resp"
)

// RegisterUser 注册用户
func RegisterUser(content *gin.Context) {
	nickname := content.PostForm("nickname")
	password := content.PostForm("password")

	userLogic := logic.UserLogic{}
	search := make(map[string]any)
	search["nickname"] = nickname
	rep := userLogic.SearchUser(search)
	if rep.Id > 0 {
		resp.Resp(resp.ReFail, "当前用户已注册", map[string]any{})
		content.Abort()
		return
	}

	user := database.User{
		Nickname: nickname,
		Password: password,
	}
	uid := user.DoRegister()
	if uid == 0 {
		resp.Resp(resp.ReSuccess, "注册失败", map[string]any{})
	}
	resp.Resp(resp.ReSuccess, "注册成功", map[string]any{"user": user})
}

// Login 用户登陆
func Login(content *gin.Context) {
	nickname := content.PostForm("nickname")
	password := content.PostForm("password")

	user := database.User{
		Nickname: nickname,
		Password: utils.Md5(password, database.UserPwdSalt),
	}

	userInfo := user.GetUserInfo()
	if userInfo.Id <= 0 {
		resp.Resp(resp.ReFail, "账号或密码错误", nil)
		content.Abort()
		return
	}
	data := make(map[string]interface{})

	tokenLogic := logic.TokenLogic{}
	j, userJwt := tokenLogic.GenerateJwt(user.Id, jwt.IndexJwtType, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = userInfo
	resp.Resp(resp.ReSuccess, "登陆成功", data)
}

func LoadUser(content *gin.Context) {
	user, ok := content.Get(string(jwt.IndexJwtType))
	if !ok {
		resp.Resp(resp.ReFail, "未查询到用户", map[string]any{})
	}
	content.JSON(200, map[string]any{"user": user})
}
