package router

import (
	"github.com/gin-gonic/gin"
	indexApi "goingo/api/index"
	"goingo/middleware"
)

type IndexRoute struct {
	group *gin.RouterGroup
}

func (r *IndexRoute) initRoute() {
	//router.GET("getAllArea", middleware.CheckJwt(), api.GetAllArea)
	r.group.GET("loadUser", middleware.CheckJwt(), indexApi.LoadUser) // 解析 jwt
	r.group.POST("registerUser", indexApi.RegisterUser)               // 注册
	r.group.POST("login", indexApi.Login)                             // 登陆
	r.group.GET("articleList", indexApi.IndexArticleList)             // 文章列表
	r.group.GET("articleDetail", indexApi.GetArticleDetail)           // 文章详情
}
