package router

import (
	"github.com/gin-gonic/gin"
	"goingo/internal/middleware"
	"goingo/internal/server/index"
)

type IndexRoute struct {
	group *gin.RouterGroup
}

func (r *IndexRoute) initRoute() {
	//router.GET("getAllArea", middleware.CheckJwt(), api.GetAllArea)
	r.group.GET("loadUser", middleware.CheckJwt(), index_api.LoadUser)                         // 解析 jwt
	r.group.POST("registerUser", index_api.RegisterUser)                                       // 注册
	r.group.POST("login", index_api.Login)                                                     // 登陆
	r.group.GET("articleList", middleware.FilterIp([]string{"*"}), index_api.IndexArticleList) // 文章列表
	r.group.GET("articleDetail", index_api.GetArticleDetail)                                   // 文章详情
	r.group.GET("getCateList", index_api.GetCateList)                                          // 分类列表
}
