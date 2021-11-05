package models

import (
	"errors"
	"gorm.io/gorm"
)

type Character struct {
	ID        int    `gorm:"column:id;not null;comment:角色ID" json:"id,omitempty"`
	Name      string `gorm:"column:name;not null;unique;comment:角色名称" json:"name,omitempty"`
	CType     uint8  `gorm:"column:admin;not null;comment:表示该角色的类型 1超级管理员，2为普通用户 超级管理员不能用户手动创建 这里可能后面会分角色等级，先预留 " json:"admin,omitempty"`
	Introduce string `gorm:"column:introduce;not null;comment:角色简介" json:"introduce,omitempty"`
}

func (u Character) TableName() string {
	return "characters"
}

// InsertCharacter 新增角色
func (c *Character) InsertCharacter(db *gorm.DB) error {
	return db.Table(c.TableName()).Create(c).Error
}

// DeleteCharacter 删除角色
func (c *Character) DeleteCharacter(db *gorm.DB) error {
	if c.ID == 0 {
		return errors.New("id is empty")
	}

	return db.Table(c.TableName()).Delete(c).Error
}

// Update 更新角色
func (c *Character) Update(db *gorm.DB) error {
	if c.ID == 0 {
		return errors.New("id is empty")
	}

	// 注意update和updates差别
	return db.Table(c.TableName()).Updates(c).Error
}
