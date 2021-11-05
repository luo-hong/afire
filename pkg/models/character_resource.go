package models

import "gorm.io/gorm"

// CharacterResource 用于描述 该角色绑定了什么资源
type CharacterResource struct {
	CID        int    `gorm:"column:c_id;primaryKey;not null;comment:角色ID" json:"uid"`
	ResourceID string `gorm:"column:res_id;primaryKey;not null;comment:资源ID" json:"res_id"`
}

func (cr CharacterResource) TableName() string {
	return "chara_res"
}

func (cr *CharacterResource) InsertCharacterResource(db *gorm.DB) error {
	return db.Table(cr.TableName()).Create(cr).Error
}

func (cr *CharacterResource) DeleteWithCid(db *gorm.DB) error {
	return db.Table(cr.TableName()).Where("c_id=?", cr.CID).Delete(CharacterResource{}).Error
}
