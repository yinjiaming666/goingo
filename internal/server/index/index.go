package index_api

import (
	"github.com/gin-gonic/gin"
	"goingo/internal/model"
	"goingo/tools/resp"
	"strconv"
)

// GetCateList 分类列表
func GetCateList(_ *gin.Context) {
	list := (new(model.Cate)).GetCateList()
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: list}).Response()
}

// IndexArticleList 前台文章列表
func IndexArticleList(c *gin.Context) {
	title := c.Query("title")
	cateId, _ := strconv.Atoi(c.Query("cate_id")) // 字符串转 int 必须要 strconv 这个包

	search := model.ArticleSearch{Title: title, Status: 0, CateId: uint(cateId)}
	article := new(model.ApiArticleList)
	list := article.GetArticleList(&search)
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: list}).Response()
}

// GetArticleDetail 文章详情
func GetArticleDetail(content *gin.Context) {
	id, _ := strconv.Atoi(content.Query("id"))
	if id <= 0 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "未查询到文章", Body: nil}).Response()
	}
	article := &model.Article{
		Id: uint(id),
	}
	article = article.GetArticleDetail()
	if article.Status == 1 {
		(&resp.JsonResp{Code: resp.ReFail, Message: "未查询到文章", Body: nil}).Response()
	}
	(&resp.JsonResp{Code: resp.ReSuccess, Message: "请求成功", Body: article}).Response()
}
