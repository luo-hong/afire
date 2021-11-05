package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	emptyCharacterResource   CharacterResource
	characterResourceField   = map[string]string{}
	characterResourcePreload = map[string]bool{}
)

type CharacterResourceSelector struct {
	PageSelector
	CID        []int
	ResourceID []string
}

func NewCharacterResourceSelector(offset, limit int) *CharacterResourceSelector {
	return &CharacterResourceSelector{
		PageSelector: MakePageSelector(offset, limit),
	}
}

func (crs *CharacterResourceSelector) makeQuery(db *gorm.DB, column ...string) (*gorm.DB, error) {

	db = db.Table(emptyCharacterResource.TableName())
	if len(column) > 0 {
		tmp, err := MakeSelectorAndPreload(db, characterResourceField, characterResourcePreload, column...)
		if err != nil {
			return nil, errors.Wrap(err, "selector tags")
		}
		db = tmp
	}

	if len(crs.CID) > 0 {
		db = db.Where("c_id in (?)", crs.CID)
	}

	if len(crs.ResourceID) > 0 {
		db = db.Where("res_id in (?)", crs.ResourceID)
	}

	return db, nil
}

func (crs *CharacterResourceSelector) Find(db *gorm.DB, column ...string) ([]CharacterResource, error) {
	var out []CharacterResource

	db, err := crs.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}

	err = crs.PageSelector.Pages(db).Find(&out).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return out, nil
}

func (crs *CharacterResourceSelector) Count(db *gorm.DB) (int64, error) {
	db, err := crs.makeQuery(db)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Count(&count).Error
	return count, err
}

func (crs *CharacterResourceSelector) Resources(db *gorm.DB) ([]string, error) {
	db, err := crs.makeQuery(db)
	if err != nil {
		return nil, err
	}

	var out []string
	err = db.Group("res_id").Pluck("res_id", &out).Error
	if err != nil {
		return nil, errors.Wrap(err, "group by")
	}

	return out, nil
}

func (crs *CharacterResourceSelector) ResourcesID(db *gorm.DB) ([]int, error) {
	db, err := crs.makeQuery(db)
	if err != nil {
		return nil, err
	}

	var out []int
	err = db.Group("c_id").Pluck("c_id", &out).Error
	if err != nil {
		return nil, errors.Wrap(err, "group by")
	}

	return out, nil
}
