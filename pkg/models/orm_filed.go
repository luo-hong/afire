package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// MakeFields 生成DB列名和结构属性名单对照表
func MakeFields(db *gorm.DB, value interface{}) (fields map[string]string, preloads map[string]bool, e error) {
	db = db.Set("_test_", "make fields")
	_ = db.Statement.Parse(value)
	fields = make(map[string]string)
	preloads = make(map[string]bool)
	for _, field := range db.Statement.Schema.Fields {
		_, foreignKey := field.TagSettings["FOREIGNKEY"]
		_, references := field.TagSettings["REFERENCES"]
		_, ignore := field.TagSettings["-"]
		if len(field.DBName) > 0 && !foreignKey && !references {
			//有列名，没外键，则是列
			fields[field.Name] = field.DBName
		} else if !ignore && len(field.DataType) == 0 {
			//无类型，则肯定为Struct，则认为它是Preloads
			preloads[field.Name] = true
		}
	}

	return fields, preloads, nil
}

// MakeSelectorAndPreload tags是struct对应的属性名，而不是DB列名
func MakeSelectorAndPreload(db *gorm.DB, fields map[string]string, preloads map[string]bool, tags ...string) (*gorm.DB, error) {
	s := make([]string, len(tags))
	if len(tags) == 0 {
		return nil, ErrEmptyField
	}
	index := 0
	tmp := db
	for _, v := range tags {
		if key, ok := fields[v]; ok {
			s[index] = key
			index++
		} else if preloads[v] {
			tmp = tmp.Preload(v)
		} else {
			return nil, errors.WithMessage(ErrNoField, v)
		}
	}

	return tmp.Select(s[:index]), nil
}
