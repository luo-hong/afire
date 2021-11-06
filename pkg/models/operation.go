package models

import (
	"reflect"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Operation struct {
	ID         uint   `gorm:"primaryKey;<-:create;comment:id"  json:"id"`
	Operator   string `gorm:"column:operator;size:16;not null;index:idx_optor;<-:create;comment:操作者" json:"operator"`
	OperatorID string `gorm:"column:operator_id;size:16;not null;index:idx_optor_id;<-:create;comment:操作者UID" json:"operator_id"`
	Operation  string `gorm:"column:operation;size:64;not null;index:idx_option;<-:create;comment:操作类型" json:"operation"`
	RequestID  string `gorm:"column:req_id;size:36;not null;index:idx_reqid;<-:create;comment:请求id" json:"req_id"`
	Details    string `gorm:"column:details;type:text;<-:create;comment:操作详情" json:"details"`
	Result     string `gorm:"column:result;type:text;comment:操作结果" json:"result"`
	UtilInfo
}

func (o Operation) TableName() string {
	return "operations"
}

// Insert 插入一个新行，将返回全新的主键ID.
func (o *Operation) Insert(db *gorm.DB) error {
	o.ID = 0 // 强制将ID置空
	if e := db.Table(o.TableName()).Create(o).Error; e != nil {
		return errors.Wrap(e, "create")
	}

	return nil
}

// Equal 忽略ID，判断其他项是否相等.
func (o *Operation) Equal(dst Operation) bool {
	dst.ID = o.ID // ID不参与比较

	return reflect.DeepEqual(o, &dst)
}
