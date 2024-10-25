package router

import (
	"app/internal/middleware"
	adminApi "app/internal/server/admin"
	"github.com/gin-gonic/gin"
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
	r.group.GET("getMenu", middleware.CheckJwt(), adminApi.GetMenu)
	r.group.POST("setMenu", middleware.CheckJwt(), adminApi.SetMenu)
	r.group.POST("delMenu", middleware.CheckJwt(), adminApi.DelMenu)
	r.group.GET("getAdminList", middleware.CheckJwt(), adminApi.GetAdminList)
}
