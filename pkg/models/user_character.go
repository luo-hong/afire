package models

import (
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

/*
	用户-角色中间表
	工号和 角色ID都是主键
*/
type UserCharacter struct {
	UID string `gorm:"column:u_id;primaryKey;not null;comment:用户工号" json:"u_id"`
	CID int    `gorm:"column:c_id;primaryKey;not null;comment:角色ID" json:"c_id"`
}

func (uc UserCharacter) TableName() string {
	return "user_character"
}

func UserCharacterBatchDelete(db *gorm.DB, uid string, cids []int) error {
	if len(cids) == 0 {
		return db.Table(emptyUserCharacter.TableName()).Where("u_id= ?", uid).Delete(emptyUserCharacter).Error

	} else {
		return db.Table(emptyUserCharacter.TableName()).Where("u_id= ? and c_id in (?)", uid, cids).Delete(emptyUserCharacter).Error
	}
}

func UserCharacterBatchInsert(db *gorm.DB, s []UserCharacter) error {
	return InsertBatches(db, emptyUserCharacter, s)
}

// Delete 删除角色下某个用户
func (c *UserCharacter) Delete(db *gorm.DB) error {

	return db.Table(c.TableName()).Delete(c).Error
}

// DeleteUserCharacter 删除uid关联角色的信息
func (c *UserCharacter) DeleteUserCharacter(db *gorm.DB) error {
	if c.UID == "" {
		return errors.New("user uid is empty")
	}
	var count int64
	if err := db.Table(c.TableName()).Where("u_id=?", c.UID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrNoRowsAffected
	}

	return db.Table(c.TableName()).Where("u_id=?", c.UID).Delete(emptyUser).Error
}
