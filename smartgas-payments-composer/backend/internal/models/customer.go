package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID               uuid.UUID `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ExternalID       string    `gorm:"external_id;type:varchar(15);not null;unique;"`
	FirstName        string    `gorm:"column:first_name;type:varchar(255);not null;default:'';"`
	FirstLastName    string    `gorm:"column:first_last_name;type:varchar(255);not null;default:'';"`
	SecondLastName   string    `gorm:"column:second_last_name;type:varchar(255);not null;default:'';"`
	PhoneNumber      string    `gorm:"column:phone_number;type:varchar(14);not null;default:'';"`
	Email            string    `gorm:"column:email;type:varchar(255);default:'';"`
	Active           bool      `gorm:"column:active;type:boolean;default:true;"`
	StripeCustomerID string    `gorm:"column:stripe_customer_id;type:varchar(255);default:'';"`
	SwitCustomerID   string    `gorm:"column:swit_customer_id;type:varchar(255);default:'';"`

	gorm.Model
}

func (c *Customer) TableName() string {
	return "customers"
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()

	return
}
