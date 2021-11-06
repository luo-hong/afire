package models

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type OperationSelector struct {
	PageSelector
	Operator   []string
	OperatorID []string
	Operation  []string
	RequestIDs []string
	StartAt    *int
	EndAt      *int
}

var (
	emptyOperation  = Operation{}
	operationFields = map[string]string{}
)

func (os *OperationSelector) makeQuery(db *gorm.DB, column ...string) (*gorm.DB, error) {
	db = db.Table(emptyOperation.TableName())
	if len(column) > 0 {
		tmp, e := MakeSelectorAndPreload(db, operationFields, nil, column...)
		if e != nil {
			return nil, errors.Wrap(e, "selector tags")
		}
		db = tmp
	}
	if len(os.RequestIDs) > 0 {
		db = db.Where("req_id in (?)", os.RequestIDs)
	}
	if len(os.OperatorID) > 0 {
		db = db.Where("operator_id in (?)", os.OperatorID)
	}
	if len(os.Operator) > 0 {
		db = db.Where("operator in (?)", os.Operator)
	}
	if len(os.Operation) > 0 {
		db = db.Where("operation in (?)", os.Operation)
	}
	if os.StartAt != nil {
		db = db.Where("created_at > ?", os.StartAt)
	}
	if os.EndAt != nil {
		db = db.Where("created_at < ?", os.EndAt)
	}

	return db, nil
}

// Find 查找操作列表.
func (os *OperationSelector) Find(db *gorm.DB, column ...string) ([]Operation, error) {
	var out []Operation
	db, e := os.makeQuery(db, column...)
	if e != nil {
		return nil, errors.Wrap(e, "make query")
	}
	e = os.PageSelector.Pages(db.Order("id DESC")).Find(&out).Error
	if e != nil && !errors.Is(e, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(e, "server select")
	}

	return out, nil
}

// Count 忽略分页值统计个数.
func (os *OperationSelector) Count(db *gorm.DB) (int, error) {
	db, e := os.makeQuery(db)
	if e != nil {
		return 0, errors.Wrap(e, "make query")
	}
	var count int64
	if e = db.Count(&count).Error; e != nil {
		return 0, errors.Wrap(e, "count")
	}

	return int(count), nil
}

// GroupWithoutCondition 忽略所有条件，执行group.
func (os *OperationSelector) GroupWithoutCondition(db *gorm.DB, column string) ([]string, error) {
	if _, ok := operationFields[column]; !ok {
		return nil, errors.Errorf("column no found: %v", column)
	}
	db, e := os.makeQuery(db)
	if e != nil {
		return nil, errors.Wrap(e, "make query")
	}

	var out []string
	e = db.Group(operationFields[column]).Pluck(operationFields[column], &out).Error
	if e != nil {
		return nil, errors.Wrap(e, "group by")
	}

	return out, nil
}
