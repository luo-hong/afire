package models

import (
	"afire/pkg/tool"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	emptyCharacter   Character
	characterField   = map[string]string{}
	characterPreload = map[string]bool{}
)

type CharacterSelector struct {
	PageSelector
	ID       []int
	NameLike string
	Admin    int
}

func NewCharacterSelector(offset, limit int) *CharacterSelector {
	return &CharacterSelector{
		PageSelector: MakePageSelector(offset, limit),
	}
}

func (cs *CharacterSelector) makeQuery(db *gorm.DB, column ...string) (*gorm.DB, error) {
	db = db.Table(emptyCharacter.TableName())
	if len(column) > 0 {
		tmp, err := MakeSelectorAndPreload(db, characterField, characterPreload, column...)
		if err != nil {
			return nil, errors.Wrap(err, "selector tags")
		}
		db = tmp
	}

	if len(cs.ID) > 0 {
		db = db.Where("id in (?)", cs.ID)
	}
	if len(cs.NameLike) > 0 {
		tmpNameLike := tool.MakeFuzzyFiled(cs.NameLike, tool.ContainsFuzzyFiled)
		db = db.Where("name like ?", tmpNameLike)
	}

	if cs.Admin > 0 {
		db = db.Where("admin = ?", cs.Admin)
	}

	return db, nil
}

// Find 查询全部角色列表
func (cs *CharacterSelector) Find(db *gorm.DB, column ...string) ([]Character, error) {
	var out []Character

	db, err := cs.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}

	err = cs.PageSelector.Pages(db).Find(&out).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return out, nil
}

func (cs *CharacterSelector) Count(db *gorm.DB) (int64, error) {
	db, err := cs.makeQuery(db)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Count(&count).Error
	return count, err
}

func (cs *CharacterSelector) AdminCharacter(db *gorm.DB) ([]int, error) {
	db, err := cs.makeQuery(db)
	if err != nil {
		return nil, err
	}

	var out []int
	err = db.Group("id").Pluck("id", &out).Error
	if err != nil {
		return nil, errors.Wrap(err, "group by")
	}

	return out, nil
}
