package model

import (
	sysLog "app/tools/logger"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var err error

const UserPwdSalt = "test"

type DbConf struct {
	UserName string
	Password string
	Ip       string
	Port     string
	DbName   string
}

func InitDb(c *DbConf) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.UserName,
		c.Password,
		c.Ip,
		c.Port,
		c.DbName,
	)

	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		// gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
			TablePrefix:   "b_", // 表前缀 https://gorm.io/zh_CN/docs/gorm_config.html
		},
	})

	if err != nil {
		sysLog.Error("连接数据库失败，请检查参数：", err)
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

}

func Db() *gorm.DB {
	return db.Debug()
}

func Paginate(p *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Page <= 0 {
			p.Page = 1
		}
		if p.PageSize <= 0 {
			p.PageSize = 20
		}
		offset := (p.Page - 1) * p.PageSize
		return Db().Offset(offset).Limit(p.PageSize)
	}
}

type Pagination struct {
	Page     int
	PageSize int
}
