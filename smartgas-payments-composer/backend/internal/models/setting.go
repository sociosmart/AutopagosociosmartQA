package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	ID    uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Name  string    `gorm:"name;type:varchar(255);not null;unique;check:name <> '';"`
	Value string    `gorm:"column:value;type:varchar(255);not null;default:'';"`

	gorm.Model
}

func (s *Setting) TableName() string {
	return "settings"
}

func (s *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()

	return
}
