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
	id, _ := strconv.Atoi(content.PostForm("id"))
	cateId, _ := strconv.Atoi(content.PostForm("cate_id"))
	articleType, _ := strconv.Atoi(content.PostForm("type"))

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

// AdminLogin 管理员登录
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
	adminId := c.GetUint(string(jwt.AdminJwtType))
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: logic.GetAdminAuth(adminId)}).Response()
}

// GetAdminList 获取管理员列表
func GetAdminList(_ *gin.Context) {
	admin := model2.Admin{}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "登陆成功", Body: admin.GetList(0)}).Response()
}

// GetRoles 获取路由
func GetRoles(c *gin.Context) {
	t, ok := c.GetQuery("type")
	var tt int
	if !ok {
		tt = -1
	} else {
		tt, _ = conv.Conv[int](t)
	}
	adminId := c.GetUint(string(jwt.AdminJwtType))
	m := logic.GetAdminAuth(adminId)
	auth := logic.GetAdminAuth(m.Id)
	menus := auth.GetAllRules(tt)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: menus}).Response()
}

// SetRoles 设置路由
func SetRoles(c *gin.Context) {
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

	menuMeta := model2.RolesMeta{}
	if v, ok := meta["title"]; ok {
		menuMeta.Title = v
	}
	if v, ok := meta["icon"]; ok {
		menuMeta.Icon = v
	}
	menu := &model2.Roles{
		Id:        id,
		Component: component,
		Meta:      menuMeta,
		Name:      name,
		Path:      path,
		Pid:       pid,
		Type:      t,
	}

	menu = menu.SetRoles()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: menu}).Response()
}

// DelRoles 删除路由
func DelRoles(c *gin.Context) {
	ids := c.PostForm("ids")
	l := strings.Split(ids, ",")
	var temp []int
	for _, v := range l {
		i, err := conv.Conv[int](v)
		if err == nil {
			temp = append(temp, i)
		}
	}
	article := &model2.Roles{}
	article.DelRoles(temp)
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
	} else {
		a := c.GetUint(string(jwt.AdminJwtType))
		aa := logic.GetAdminAuth(a)
		admin.Pid = aa.Id
	}

	admin.SetAdmin()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "更新成功", Body: admin}).Response()
}

// DelAdmin 删除管理员
func DelAdmin(c *gin.Context) {
	ids := c.PostForm("ids")
	l := strings.Split(ids, ",")
	var temp []int
	for _, v := range l {
		i, err := conv.Conv[int](v)
		if err == nil {
			temp = append(temp, i)
		}
	}
	article := &model2.Admin{}
	article.DelAdmin(temp)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "成功", Body: nil}).Response()
}

// GetRolesGroupList 获取角色组列表
func GetRolesGroupList(c *gin.Context) {
	adminId := c.GetUint(string(jwt.AdminJwtType))
	if logic.GetAdminAuth(adminId).IsSuperAdmin {
		adminId = 0
	}
	list := (new(model2.RolesGroup)).GetRolesGroupList(adminId)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: list}).Response()
}

// SetRolesGroup 设置角色组
func SetRolesGroup(c *gin.Context) {
	rolesIds := conv.ConvPostForm[string](c, "roles_ids")
	id := conv.ConvPostForm[uint](c, "id")
	name := conv.ConvPostForm[string](c, "name")
	adminId := c.GetUint(string(jwt.AdminJwtType))

	group := model2.RolesGroup{}
	if id > 0 {
		group.Id = id
		// 校验是否可以修改该角色组
		group.GetRolesGroup()
		if group.AdminId != adminId {
			(&resp.JsonResp{Code: resp.ReFail, Message: "无权限操作此角色组", Body: nil}).Response()
		}
	}

	if name != "" {
		group.Name = name
	}
	if rolesIds == "" {
		(&resp.JsonResp{Code: resp.ReFail, Message: "roles_ids 不能为空", Body: nil}).Response()
	}
	rolesIdsList, err := conv.Explode[uint](",", rolesIds)
	if err != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "roles_ids 格式错误", Body: nil}).Response()
	}
	roles := model2.Roles{}
	getRoles := roles.GetRoles(rolesIdsList, -1)
	if len(getRoles) != len(rolesIdsList) {
		(&resp.JsonResp{Code: resp.ReFail, Message: "illegal roles_ids", Body: nil}).Response()
	}

	a := logic.GetAdminAuth(adminId)
	if a.AuthRules(rolesIdsList) != nil {
		(&resp.JsonResp{Code: resp.ReFail, Message: "无权限操作此模块", Body: nil}).Response()
	}
	group.AdminId = adminId
	rolesGroup := group.SetRolesGroup()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "操作成功", Body: rolesGroup}).Response()
}

// DelRolesGroup 删除角色组
func DelRolesGroup(c *gin.Context) {
	groupId := conv.ConvPostForm[uint](c, "group_id")
	if groupId < 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "缺少参数", Body: nil}).Response()
	}
	adminId := c.GetUint(string(jwt.AdminJwtType))
	group := model2.RolesGroup{}
	group.Id = groupId
	group.GetRolesGroup()
	if group.AdminId != adminId {
		(&resp.JsonResp{Code: resp.ReFail, Message: "无权限操作此角色组", Body: nil}).Response()
	}
	group.DelRolesGroup()
	adminGroup := model2.AdminRolesGroup{}
	adminGroup.DelByRolesGroupId(groupId)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "删除成功", Body: nil}).Response()
}
