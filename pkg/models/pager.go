package models

import "gorm.io/gorm"

type PageSelector struct {
	Offset *int
	Limit  *int
}

func MakePageSelector(offset, limit int) PageSelector {
	var addrOffset *int
	if offset > 0 {
		addrOffset = &offset
	}
	var addrLimit *int
	if limit >= 0 {
		addrLimit = &limit
	}

	return PageSelector{
		Offset: addrOffset,
		Limit:  addrLimit,
	}
}

// Pages 对db执行分页操作
//     selector.Pages(db).Query(&out)
func (p PageSelector) Pages(db *gorm.DB) *gorm.DB {
	if p.Offset != nil {
		db = db.Offset(*p.Offset)
	}
	if p.Limit != nil {
		db = db.Limit(*p.Limit)
	}

	return db
}
