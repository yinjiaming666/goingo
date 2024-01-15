package database

import (
	"goingo/logger"
	"goingo/utils/conv"
	"reflect"
)

type BaseModel interface {
	IsBaseModel()
	ToMap(child BaseModel) map[string]any
	InitWithMap(a map[string]any, child BaseModel)
	CreateTable(child BaseModel)
}

type MysqlBaseModel struct{}

func (m *MysqlBaseModel) IsBaseModel() {}

func (m *MysqlBaseModel) ToMap(child BaseModel) map[string]any {
	return conv.Struct2Map(child, true)
}

func (m *MysqlBaseModel) InitWithMap(arr map[string]any, child BaseModel) {
	conv.Map2Struct(arr, child, true)
}

// CreateTable 生成数据表
func (m *MysqlBaseModel) CreateTable(child BaseModel) {
	v := reflect.ValueOf(child)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	logger.Info("init mysql table " + v.String())
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(child)
	if err != nil {
		logger.Error("数据表生成失败", err.Error())
		return
	}
}
