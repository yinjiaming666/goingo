package admin_api

import (
	"app/internal/logic"
	model2 "app/internal/model"
	"app/tools"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

// SetArticle 添加修改文章
func SetArticle(content *gin.Context) {
	title := content.PostForm("title")
	contents := content.PostForm("content")
	status, _ := strconv.Atoi(content.PostForm("status"))
	id, _ := strconv.Atoi(content.PostForm("id"))            // 字符串转 int 必须要 strconv 这个包
	cateId, _ := strconv.Atoi(content.PostForm("cate_id"))   // 字符串转 int 必须要 strconv 这个包
	articleType, _ := strconv.Atoi(content.PostForm("type")) // 字符串转 int 必须要 strconv 这个包

	article := &model2.Article{
		Id:      uint(id),
		Title:   title,
		Content: contents,
		Status:  int8(status),
		CateId:  cateId,
		Type:    uint(articleType),
	}

	article = article.SetArticle()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: nil}).Response()
}

// DelArticle 删除文章
func DelArticle(content *gin.Context) {
	id, _ := strconv.Atoi(content.PostForm("id"))
	article := &model2.Article{
		Id: uint(id),
	}
	article.DelArticle()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: nil}).Response()
}

// ArticleList 后台文章列表
func ArticleList(content *gin.Context) {
	title := content.Query("title")

	var status int
	if content.Query("status") == "" {
		status = 99
	} else {
		status, _ = strconv.Atoi(content.Query("status"))
	}

	var cateId int
	if content.Query("cate_id") == "" {
		cateId = 99
	} else {
		cateId, _ = strconv.Atoi(content.Query("cate_id"))
	}

	uid, _ := strconv.Atoi(content.Query("uid"))
	search := &model2.ArticleSearch{
		Title:  title,
		Status: int8(status),
		Uid:    uint(uid),
		CateId: uint(cateId),
	}

	var article model2.ApiArticleList
	data := article.GetArticleList(search)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: data}).Response()
}

// AdminLogin Login 管理员登录
func AdminLogin(content *gin.Context) {
	account := content.PostForm("account")
	password := content.PostForm("password")

	if account == "" || password == "" {
		(&resp.JsonResp{Code: resp.ReFail, Message: "请输入账号密码", Body: nil}).Response()
	}

	admin := &model2.Admin{
		Account:  account,
		Password: tools.Md5(password, model2.UserPwdSalt),
	}

	admin = admin.GetAdmin()
	if admin.Id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "账号密码错误", Body: nil}).Response()
	}

	data := make(map[string]interface{})
	tokenLogic := logic.TokenLogic{}
	j, userJwt := tokenLogic.GenerateJwt(admin.Id, jwt.AdminJwtType, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = admin
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: data}).Response()
}

// GetAdminInfo 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	admin, _ := c.Get(string(jwt.AdminJwtType))
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: admin}).Response()
}

func GetMenu(c *gin.Context) {
	t, ok := c.GetQuery("type")
	var tt int
	if !ok {
		tt = -1
	} else {
		tt, _ = conv.Conv[int](t)
	}
	menu := model2.Role{}
	menus := menu.GetMenusByRoleIds("*", tt)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: menus}).Response()
}

// GetCateList 获取分类列表
func GetCateList(_ *gin.Context) {
	list := (new(model2.Cate)).GetCateList()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: list}).Response()
}

// SetAdminInfo 更新管理员信息
func SetAdminInfo(c *gin.Context) {
	temp, _ := c.Get(string(jwt.AdminJwtType))
	admin, ok := temp.(*model2.Admin)
	if !ok {
		(&resp.JsonResp{Code: resp.ReFail, Message: "未查询到账号", Body: nil}).Response()
	}

	name := c.PostForm("name")
	password := c.PostForm("password")
	avatar := c.PostForm("avatar")

	data := make(map[string]interface{})
	if name != "" {
		data["name"] = name
	}

	if avatar != "" {
		data["avatar"] = avatar
	}

	if password != "" {
		data["password"] = tools.Md5(password, model2.UserPwdSalt)
	}

	admin = &model2.Admin{Id: admin.Id}
	admin = admin.UpdateAdmin(data)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "更新成功", Body: admin}).Response()
}

// GetArticleDetail 文章详情
func GetArticleDetail(content *gin.Context) {
	id, _ := strconv.Atoi(content.Query("id"))
	if id <= 0 {
		(&resp.JsonResp{Code: resp.ReSuccess, Message: "未查询到文章", Body: model2.Article{}}).Response()
	}
	article := &model2.Article{
		Id: uint(id),
	}
	article = article.GetArticleDetail()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: article}).Response()
}
