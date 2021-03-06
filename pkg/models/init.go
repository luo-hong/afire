package models

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
)

func createDB(db *gorm.DB, table interface{}) error {
	return db.AutoMigrate(table)
}

func InitModels(db *gorm.DB) error {
	var e error

	// 创建user表 用户信息表
	fmt.Println("[models] user")
	userField, userPreload, e = MakeFields(db, &User{}) // 后续匹配数据用到的参数
	if e != nil {
		return errors.Wrap(e, "auto fields users")
	}
	e = createDB(db, &User{})
	if e != nil {
		return errors.Wrap(e, "auto migrate users")
	}

	// 创建user_character表 用户匹配表
	fmt.Println("[models] user character")
	userCharacterField, userCharacterPreload, e = MakeFields(db, &UserCharacter{})
	if e != nil {
		return errors.Wrap(e, "auto fields user_character")
	}
	e = createDB(db, &UserCharacter{})
	if e != nil {
		return errors.Wrap(e, "auto migrate user_character")
	}

	// 创建character表  角色信息表
	fmt.Println("[models] character")
	characterField, characterPreload, e = MakeFields(db, &Character{})
	if e != nil {
		return errors.Wrap(e, "auto fields character")
	}
	e = createDB(db, &Character{})
	if e != nil {
		return errors.Wrap(e, "auto migrate character")
	}

	// 创建character res表 用户的资源表
	fmt.Println("[models] character res")
	characterResourceField, characterResourcePreload, e = MakeFields(db, &CharacterResource{})
	if e != nil {
		return errors.Wrap(e, "auto fields character_resource")
	}
	e = createDB(db, &CharacterResource{})
	if e != nil {
		return errors.Wrap(e, "auto migrate character_resource")
	}

	// 创建operation表  记录用户操作详情
	fmt.Println("[models] operations")
	operationFields, _, e = MakeFields(db, &Operation{})
	if e != nil {
		return errors.Wrap(e, "make fields operation")
	}
	e = createDB(db, &Operation{})
	if e != nil {
		return errors.Wrap(e, "auto migrate operation")
	}

	return nil
}

type TableNameInit interface {
	TableName() string
}

func InsertBatches(db *gorm.DB, tn TableNameInit, f interface{}) error {

	v := reflect.ValueOf(f)

	if v.Kind() != reflect.Slice {
		return errors.New("type is illegal")
	}

	if v.Len() < 1 {
		return nil
	}

	return db.Table(tn.TableName()).CreateInBatches(f, v.Len()).Error
}
