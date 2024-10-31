package admin_api

import (
	"app/internal/logic"
	model2 "app/internal/model"
	"app/tools"
	"app/tools/conv"
	"app/tools/jwt"
	"app/tools/resp"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
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

// GetAdminList 获取管理员列表
func GetAdminList(c *gin.Context) {
	admin := model2.Admin{}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: admin.GetList(0)}).Response()
}

// GetMenu 获取路由
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

// SetMenu 设置路由
func SetMenu(c *gin.Context) {
	id, err := conv.Conv[uint](c.PostForm("id"))
	if err != nil {
		fmt.Println(err)
	}
	pid, err := conv.Conv[uint](c.PostForm("pid"))
	if err != nil {
		fmt.Println(err)
	}
	t, err := conv.Conv[uint](c.PostForm("type"))
	if err != nil {
		fmt.Println(err)
	}

	name := c.PostForm("name")
	path := c.PostForm("path")
	component := c.PostForm("component")
	meta := c.PostFormMap("meta")

	menuMeta := model2.MenuMeta{}
	if v, ok := meta["title"]; ok {
		menuMeta.Title = v
	}
	if v, ok := meta["icon"]; ok {
		menuMeta.Icon = v
	}
	menu := &model2.Menu{
		Id:        id,
		Component: component,
		Meta:      menuMeta,
		Name:      name,
		Path:      path,
		Pid:       pid,
		Type:      t,
	}

	menu = menu.SetMenu()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: menu}).Response()
}

// DelMenu 删除路由
func DelMenu(c *gin.Context) {
	ids := c.PostForm("ids")
	l := strings.Split(ids, ",")
	var temp []int
	for _, v := range l {
		i, err := conv.Conv[int](v)
		if err == nil {
			temp = append(temp, i)
		}
	}
	article := &model2.Menu{}
	article.DelMenu(temp)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: nil}).Response()
}

// GetCateList 获取分类列表
func GetCateList(_ *gin.Context) {
	list := (new(model2.Cate)).GetCateList()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: list}).Response()
}

// SetAdminInfo 更新管理员信息
func SetAdminInfo(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	avatar := c.PostForm("avatar")
	id := c.PostForm("id")
	admin := model2.Admin{}

	if name != "" {
		admin.Name = name
	}

	if avatar != "" {
		admin.Avatar = avatar
	}

	if password != "" {
		admin.Password = tools.Md5(password, model2.UserPwdSalt)
	}

	if id != "" {
		i, _ := conv.Conv[uint](id)
		admin.Id = i
	}

	admin.SetAdmin()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "更新成功", Body: admin}).Response()
}
