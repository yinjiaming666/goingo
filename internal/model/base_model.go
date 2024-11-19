package model

import (
	"app/tools/conv"
	"app/tools/logger"
	"fmt"
	"reflect"
)

type BaseModel interface {
	IsModel()
	ToMap(child BaseModel) map[string]any
	InitWithMap(a map[string]any, child BaseModel)
	CreateTable(child BaseModel)
	GetTableComment() string
	SetTableComment(comment string) BaseModel
	GetEngin() string
	SetEngin(engin string) BaseModel
}

type MysqlBaseModel struct {
	TableComment string // 表注释
	Engin        string // 表注释
}

func (m *MysqlBaseModel) IsModel() {}

func (m *MysqlBaseModel) GetTableComment() string {
	return m.TableComment
}

func (m *MysqlBaseModel) SetTableComment(comment string) BaseModel {
	m.TableComment = comment
	return m
}

func (m *MysqlBaseModel) GetEngin() string {
	if m.Engin == "" {
		return "InnoDB"
	}
	return m.Engin
}

func (m *MysqlBaseModel) SetEngin(engin string) BaseModel {
	m.Engin = engin
	return m
}

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
	logger.Debug("INIT MYSQL TABLE " + v.String())
	err := db.Set("gorm:table_options", fmt.Sprintf("ENGINE=%s, COMMENT='%s'", m.GetEngin(), m.GetTableComment())).AutoMigrate(child)
	if err != nil {
		logger.Error("数据表生成失败", err.Error())
		return
	}
}
