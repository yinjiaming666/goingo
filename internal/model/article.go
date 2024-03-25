package model

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

type Article struct {
	*MysqlBaseModel `gorm:"-:all"`        // -:all 无读写迁移权限，该字段不在数据库中
	Id              uint                  `json:"id" gorm:"primaryKey;type:INT(8) UNSIGNED NOT NULL AUTO_INCREMENT"`
	Title           string                `json:"title" gorm:"type:VARCHAR(1200) NOT NULL;default:''"`
	Content         string                `json:"content" gorm:"type:TEXT;"`
	Uid             uint                  `json:"uid" gorm:"type:INT UNSIGNED NOT NULL;default:0"`
	Status          int8                  `json:"status" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
	ViewNum         uint                  `json:"view_num" gorm:"type:INT UNSIGNED NOT NULL;default:0"`
	CateId          int                   `json:"cate_id" gorm:"type:INT UNSIGNED NOT NULL;default:0"`
	CreateTime      int64                 `json:"create_time" gorm:"autoCreateTime;type:BIGINT UNSIGNED NOT NULL;"` // 自动写入时间戳
	CreateTimeStr   string                `json:"create_time_str" gorm:"-:all"`                                     // -:all 无读写迁移权限，该字段不在数据库中
	DeleteTime      soft_delete.DeletedAt `json:"delete_time" gorm:"type:BIGINT UNSIGNED NOT NULL;default:0"`
	Cate            Cate                  `json:"cate" gorm:"foreignKey:cate_id;references:id;-:migration"` // -:migration 表示无迁移权限
	Type            uint                  `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`

	// belongsTo 文章分类(article 的 cate_id 指向 cate 的 id)  查询的时候使用 Joins 或者 Preload
	// 通过 debug 推测 Joins 是 leftJoin ，Preload 应该是循环查询的
}

type ApiArticleList struct {
	Id            uint   `json:"id" gorm:"primaryKey"`
	Title         string `json:"title"`
	Status        int8   `json:"status"`
	ViewNum       uint   `json:"view_num"`
	CreateTime    int64  `json:"create_time" gorm:"autoCreateTime"`
	CreateTimeStr string `json:"create_time_str" gorm:"-:all"`
	CateId        uint   `json:"cate_id" gorm:"default:1"` // 默认值为 1
	Cate          Cate   `json:"cate" gorm:"foreignKey:cate_id;references:id"`
	Type          uint   `json:"type" gorm:"type:TINYINT(8) UNSIGNED NOT NULL;default:0"`
}

type ArticleSearch struct {
	Uid    uint
	Status int8
	Title  string
	CateId uint
}

func (article *Article) SetArticle() *Article {
	if article.Id <= 0 {
		db.Debug().Create(&article)
	} else {
		db.Debug().Select("title", "content", "status", "cate_id", "type").Model(&article).Updates(&article)
	}
	return article
}

func (article *Article) DelArticle() *Article {
	db.Delete(&article, article.Id)
	return article
}

func (article *ApiArticleList) GetArticleList(search *ArticleSearch) []ApiArticleList {
	list := make([]ApiArticleList, 0)
	where := make(map[string]interface{})

	if search.Uid > 0 {
		where["uid"] = search.Uid
	}
	if search.Status != 99 {
		where["status"] = search.Status
	}
	if search.CateId != 99 && search.CateId != 0 {
		where["cate_id"] = search.CateId
	}

	conn := db.Debug().Where(where)

	if search.Title != "" {
		conn.Where("b_article.title LIKE ?", "%"+search.Title+"%") // 链式操作
	}

	a := new(Article)
	conn.Model(a).Joins("Cate").Order("id desc").Find(&list)
	return list
}

func (article *Article) GetArticleDetail() *Article {
	db.Where(article).Debug().First(article)
	return article
}

// AfterFind 查询钩子，在查询后会执行
func (article *Article) AfterFind(*gorm.DB) (err error) {
	if article.CreateTime != 0 {
		article.CreateTimeStr = time.Unix(article.CreateTime, 0).Format("2006-01-02 15:04:05")
	}
	return
}

// AfterFind 查询钩子，在查询后会执行
func (article *ApiArticleList) AfterFind(*gorm.DB) (err error) {
	if article.CreateTime != 0 {
		article.CreateTimeStr = time.Unix(article.CreateTime, 0).Format("2006-01-02 15:04:05")
	}
	return
}
