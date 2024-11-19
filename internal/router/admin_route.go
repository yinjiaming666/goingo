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
	r.group.GET("getRoles", middleware.CheckJwt(), adminApi.GetRoles)
	r.group.POST("setRoles", middleware.CheckJwt(), adminApi.SetRoles)
	r.group.POST("delRoles", middleware.CheckJwt(), adminApi.DelRoles)
	r.group.GET("getAdminList", middleware.CheckJwt(), adminApi.GetAdminList)
	r.group.GET("delAdmin", middleware.CheckJwt(), adminApi.DelAdmin)
	r.group.GET("setRolesGroup", middleware.CheckJwt(), adminApi.SetRolesGroup)
	r.group.GET("getRolesGroupList", middleware.CheckJwt(), adminApi.GetRolesGroupList)
	r.group.GET("delRolesGroup", middleware.CheckJwt(), adminApi.DelRolesGroup)
}
