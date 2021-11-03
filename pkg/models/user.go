package models

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
)

type User struct {
	UID   string `gorm:"column:uid;<-:create;not null;primaryKey;comment:工号" json:"id"`
	Name  string `gorm:"column:name;not null;comment:用户名称" json:"name"`
	Pwd   string `gorm:"column:pwd;size:64;not null;comment:密码"`
	Phone string `gorm:"column:phone;not null;comment:手机号"`
	Email string `gorm:"column:email;comment:邮箱"`
}

func (m User) TableName() string {
	return "users"
}

func (m User) Insert(db *gorm.DB) error {
	return db.Table(m.TableName()).Create(m).Error
}

func (m User) Update(db *gorm.DB) error {
	if m.UID == "" {
		return errors.New("user uid is empty")
	}
	m.Pwd = ""
	var count int64
	if err := db.Table(m.TableName()).Where("uid=?", m.UID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRowsAffected
	}

	return db.Table(m.TableName()).Omit("pwd").Updates(m).Error
}

func (m User) UpdateByColumn(db *gorm.DB, columns ...string) error {
	if m.UID == "" {
		return errors.New("user uid is empty")
	}
	d := reflect.ValueOf(m)
	upMap := map[string]interface{}{}
	tmpVal := reflect.Indirect(d)
	for _, v := range columns {
		c, ok := userField[v]
		if ok {
			upMap[c] = tmpVal.FieldByName(v).Interface()
		}
	}

	return db.Table(m.TableName()).Where("uid=?", m.UID).Updates(upMap).Error
}

func (m User) Delete(db *gorm.DB) error {
	if m.UID == "" {
		return errors.New("user uid is empty")
	}
	var count int64
	if err := db.Table(m.TableName()).Where("uid=?", m.UID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRowsAffected
	}

	return db.Table(m.TableName()).Where("uid=?", m.UID).Delete(emptyUser).Error
}
