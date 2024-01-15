package router

import (
	"github.com/gin-gonic/gin"
	"goingo/internal/middleware"
	adminApi "goingo/internal/server/admin"
)

type AdminRoute struct {
	group *gin.RouterGroup
}

func (r *AdminRoute) initRoute() {
	r.group.POST("login", adminApi.AdminLogin)
	r.group.GET("getArticleDetail", adminApi.GetArticleDetail)
	r.group.POST("setArticle", middleware.CheckJwt(), adminApi.SetArticle)
	r.group.POST("delArticle", middleware.CheckJwt(), adminApi.DelArticle)
	r.group.GET("articleList", middleware.CheckJwt(), adminApi.ArticleList)
	r.group.GET("getAdminInfo", middleware.CheckJwt(), adminApi.GetAdminInfo)
	r.group.GET("getCateList", middleware.CheckJwt(), adminApi.GetCateList)
	r.group.POST("setAdminInfo", middleware.CheckJwt(), adminApi.SetAdminInfo)
}
