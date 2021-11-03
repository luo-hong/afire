package models

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func createDB(db *gorm.DB, table interface{}) error {
	return db.AutoMigrate(table)
}

func InitModels(db *gorm.DB) error {
	var e error

	// 创建user表
	fmt.Println("[models] user")
	//userField, userPreload, e = MakeFields(db, &User{}) // 后续匹配数据用到的参数
	if e != nil {
		return errors.Wrap(e, "auto fields users")
	}
	e = createDB(db, &User{})
	if e != nil {
		return errors.Wrap(e, "auto migrate users")
	}

	return nil
}
