package index_api

import (
	"github.com/gin-gonic/gin"
	"goingo/database"
	"goingo/utils/resp"
	"strconv"
)

// IndexArticleList 前台文章列表
func IndexArticleList(c *gin.Context) {
	title := c.Query("title")
	cateId, _ := strconv.Atoi(c.PostForm("cate_id")) // 字符串转 int 必须要 strconv 这个包

	search := database.ArticleSearch{Title: title, Status: 0, CateId: uint(cateId)}
	article := new(database.ApiArticleList)
	list := article.GetArticleList(&search)

	resp.Resp(resp.ReSuccess, "请求成功", list)
}

// GetArticleDetail 文章详情
func GetArticleDetail(content *gin.Context) {
	id, _ := strconv.Atoi(content.Query("id"))
	if id <= 0 {
		resp.Resp(resp.ReSuccess, "未查询到文章", nil)
	}
	article := &database.Article{
		Id: uint(id),
	}
	article = article.GetArticleDetail()
	if article.Status == 1 {
		resp.Resp(resp.ReSuccess, "未查询到文章", nil)
	}
	resp.Resp(resp.ReSuccess, "请求成功", article)
}
