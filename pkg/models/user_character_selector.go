package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	emptyUserCharacter   UserCharacter
	userCharacterField   = map[string]string{}
	userCharacterPreload = map[string]bool{}
)

type UserCharacterSelector struct {
	PageSelector
	UID []string
	CID []int
}

func NewUserCharacterSelector(offset, limit int) *UserCharacterSelector {
	return &UserCharacterSelector{
		PageSelector: MakePageSelector(offset, limit),
	}
}

func (ucs *UserCharacterSelector) makeQuery(db *gorm.DB, column ...string) (*gorm.DB, error) {
	db = db.Table(emptyUserCharacter.TableName())
	if len(column) > 0 {
		tmp, err := MakeSelectorAndPreload(db, userCharacterField, userCharacterPreload, column...)
		if err != nil {
			return nil, errors.Wrap(err, "selector tags")
		}
		db = tmp
	}

	if len(ucs.UID) > 0 {
		db = db.Where("u_id in (?)", ucs.UID)
	}

	if len(ucs.CID) > 0 {
		db = db.Where("c_id in (?)", ucs.CID)
	}

	return db, nil
}

func (ucs *UserCharacterSelector) Find(db *gorm.DB, column ...string) ([]UserCharacter, error) {
	var out []UserCharacter

	db, err := ucs.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}

	err = ucs.PageSelector.Pages(db).Find(&out).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return out, nil
}

func (ucs *UserCharacterSelector) Count(db *gorm.DB) (int64, error) {
	db, err := ucs.makeQuery(db)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Count(&count).Error
	return count, err
}

func (ucs *UserCharacterSelector) UIDs(db *gorm.DB) ([]string, error) {
	db, e := ucs.makeQuery(db)
	if e != nil {
		return nil, e
	}
	var out []string
	e = db.Group("u_id").Pluck("u_id", &out).Error
	if e != nil {
		return nil, errors.Wrap(e, "group by")
	}

	return out, nil
}

func (ucs *UserCharacterSelector) CIDs(db *gorm.DB) ([]int, error) {
	db, e := ucs.makeQuery(db)
	if e != nil {
		return nil, e
	}
	var out []int
	e = db.Group("c_id").Pluck("c_id", &out).Error
	if e != nil {
		return nil, errors.Wrap(e, "group by")
	}

	return out, nil
}
