package models

type UtilInfo struct {
	Editor    string `gorm:"column:editor;not null;size:16;comment:创建人" json:"editor,omitempty"`
	CreatedAt int    `gorm:"index;<-:create;not null;comment:创建时间" json:"created_at,omitempty"`
	UpdatedAt int    `gorm:"not null;comment:最后更新时间" json:"updated_at,omitempty"`
}
