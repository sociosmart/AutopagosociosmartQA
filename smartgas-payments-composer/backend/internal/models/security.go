package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorizedApplication struct {
	ID              uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ApplicationName string    `gorm:"column:application_name;type:varchar(100);not null;default:'';"`
	AppKey          uuid.UUID `gorm:"column:app_key;type:varchar(36);not null;unique;"`
	ApiKey          uuid.UUID `gorm:"column:api_key;type:varchar(36);not null;unique;"`
	Active          bool      `gorm:"column:active;type:boolean;not null;default:true;"`
	gorm.Model
}

func (ap *AuthorizedApplication) TableName() string {
	return "authorized_applications"
}

func (ap *AuthorizedApplication) BeforeCreate(tx *gorm.DB) (err error) {
	ap.ID = uuid.New()

	ap.ApiKey = uuid.New()
	ap.AppKey = uuid.New()

	return
}
