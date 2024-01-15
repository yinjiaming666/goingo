package admin_api

import (
	"github.com/gin-gonic/gin"
	"goingo/internal/logic"
	model2 "goingo/internal/model"
	"goingo/tools"
	"goingo/tools/jwt"
	"goingo/tools/resp"
	"strconv"
)

// SetArticle 添加修改文章
func SetArticle(content *gin.Context) {
	title := content.PostForm("title")
	contents := content.PostForm("content")
	status, _ := strconv.Atoi(content.PostForm("status"))
	id, _ := strconv.Atoi(content.PostForm("id"))          // 字符串转 int 必须要 strconv 这个包
	cateId, _ := strconv.Atoi(content.PostForm("cate_id")) // 字符串转 int 必须要 strconv 这个包

	article := &model2.Article{
		Id:      uint(id),
		Title:   title,
		Content: contents,
		Status:  int8(status),
		CateId:  cateId,
	}

	article = article.SetArticle()
	resp.Resp(resp.ReSuccess, "成功", nil)
}

// DelArticle 删除文章
func DelArticle(content *gin.Context) {
	id, _ := strconv.Atoi(content.PostForm("id"))
	article := &model2.Article{
		Id: uint(id),
	}
	article.DelArticle()
	resp.Resp(resp.ReSuccess, "成功", nil)
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
	resp.Resp(resp.ReSuccess, "成功", data)
}

// AdminLogin Login 管理员登录
func AdminLogin(content *gin.Context) {
	account := content.PostForm("account")
	password := content.PostForm("password")

	if account == "" || password == "" {
		resp.Resp(resp.ReFail, "请输入账号密码", nil)
	}

	admin := &model2.Admin{
		Account:  account,
		Password: tools.Md5(password, model2.UserPwdSalt),
	}

	admin = admin.GetAdmin()
	if admin.Id <= 0 {
		resp.Resp(resp.ReFail, "账号密码错误", nil)
	}

	data := make(map[string]interface{})
	tokenLogic := logic.TokenLogic{}
	j, userJwt := tokenLogic.GenerateJwt(admin.Id, jwt.AdminJwtType, 0)
	userJwt.Token = ""
	data["token"] = j
	data["token_info"] = userJwt
	data["user"] = admin
	resp.Resp(resp.ReSuccess, "登陆成功", data)
}

// GetAdminInfo 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	admin, _ := c.Get(string(jwt.AdminJwtType))

	data := make(map[string]interface{})
	data["user"] = admin
	resp.Resp(resp.ReSuccess, "登陆成功", data)
}

// GetCateList 获取分类列表
func GetCateList(_ *gin.Context) {
	list := (new(model2.Cate)).GetCateList()
	resp.Resp(resp.ReSuccess, "请求成功", list)
}

// SetAdminInfo 更新管理员信息
func SetAdminInfo(c *gin.Context) {
	temp, _ := c.Get(string(jwt.AdminJwtType))
	admin, ok := temp.(*model2.Admin)
	if !ok {
		resp.Resp(resp.ReFail, "未查询到账号", nil)
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
	resp.Resp(resp.ReSuccess, "更新成功", admin)
}

// GetArticleDetail 文章详情
func GetArticleDetail(content *gin.Context) {
	id, _ := strconv.Atoi(content.Query("id"))
	if id <= 0 {
		resp.Resp(resp.ReSuccess, "未查询到文章", model2.Article{})
	}
	article := &model2.Article{
		Id: uint(id),
	}
	article = article.GetArticleDetail()
	resp.Resp(resp.ReSuccess, "请求成功", article)
}
