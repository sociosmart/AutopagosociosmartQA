package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID     uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Name   string    `gorm:"column:name;type:varchar(100);not null;unique;"`
	Groups []*Group  `gorm:"many2many:group_permissions;"`
	Users  []*User   `gorm:"many2many:user_permissions;"`

	gorm.Model
}

func (p *Permission) TableName() string {
	return "permissions"
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()

	return
}

type Group struct {
	ID          uuid.UUID     `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Name        string        `gorm:"column:name;type:varchar(100);not null;unique;"`
	Permissions []*Permission `gorm:"many2many:group_permissions;"`
	Users       []*User       `gorm:"many2many:user_groups;"`

	gorm.Model
}

func (g *Group) TableName() string {
	return "groups"
}

func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()

	return
}
