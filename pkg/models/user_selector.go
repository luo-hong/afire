package models

import (
	"afire/pkg/tool"
	"afire/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	emptyUser   User
	userField   = map[string]string{}
	userPreload = map[string]bool{}
)

type UserSelector struct {
	PageSelector
	UID       []string // 精准匹配
	UIDLike   string   // 模糊
	NameLike  string   // 模糊
	UIDOrName string   // 模糊
}

func (us *UserSelector) makeQuery(db *gorm.DB, column ...string) (*gorm.DB, error) {
	db = db.Table(emptyUser.TableName())
	if len(column) > 0 {
		tmp, err := MakeSelectorAndPreload(db, userField, userPreload, column...)
		if err != nil {
			return nil, errors.Wrap(err, "selector tags")
		}
		db = tmp
	}

	uids := utils.RemoveRepeatString(us.UID)
	if len(uids) > 0 {
		db = db.Where("uid in (?)", uids)
	}

	if len(us.UIDLike) > 0 {
		tmpUIDLike := tool.MakeFuzzyFiled(us.UIDLike, tool.ContainsFuzzyFiled)
		db = db.Where("uid like ?", tmpUIDLike)
	}

	if len(us.NameLike) > 0 {
		tmpNameLike := tool.MakeFuzzyFiled(us.NameLike, tool.ContainsFuzzyFiled)
		db = db.Where("uid like ?", tmpNameLike)
	}

	if len(us.UIDOrName) > 0 {
		tmpUIDOrName := tool.MakeFuzzyFiled(us.UIDOrName, tool.ContainsFuzzyFiled)
		db = db.Where("uid like ? or name like ?", tmpUIDOrName, tmpUIDOrName)
	}

	return db, nil
}
func (us *UserSelector) Find(db *gorm.DB, column ...string) ([]User, error) {
	var out []User

	db, err := us.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}

	err = us.PageSelector.Pages(db).Find(&out).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return out, nil
}

func (us *UserSelector) Count(db *gorm.DB) (int64, error) {
	db, err := us.makeQuery(db)
	if err != nil {
		return 0, err
	}
	var count int64
	err = db.Count(&count).Error
	return count, err
}

func (s *UserSelector) One(db *gorm.DB, column ...string) (*User, error) {
	db, err := s.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}
	var out User
	err = s.Pages(db).First(&out).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "select")
	}
	return &out, nil
}

func (s *UserSelector) UIDs(db *gorm.DB, column ...string) ([]string, error) {
	db, err := s.makeQuery(db, column...)
	if err != nil {
		return nil, err
	}
	var out []string
	err = db.Group("uid").Pluck("uid", &out).Error
	if err != nil {
		return nil, errors.Wrap(err, "pluck")
	}

	return out, nil
}
