package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Level struct {
	ID          uuid.UUID  `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Name        *string    `gorm:"column:name;type:varchar(100);not null;unique;"`
	Discount    *float64   `gorm:"column:discount;type:float;not null;default:0;check:discount > -1;"`
	MinCharges  *int       `gorm:"column:min_charges;type:SMALLINT;not null;default:0;check:min_charges > -1;"`
	MinAmount   *float64   `gorm:"column:min_amount;type:FLOAT;not null;default:0;check:min_amount > -1;"`
	Active      *bool      `gorm:"column:active;type:boolean;not null;default:true;"`
	CreatedByID *uuid.UUID `gorm:"column:created_by_id;type:varchar(36);"`
	CreatedBy   *User      `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
	UpdatedByID *uuid.UUID `gorm:"column:updated_by_id;type:varchar(36);"`
	UpdatedBy   *User      `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:SET NULL;"`

	gorm.Model
}

type CustomerLevel struct {
	ID              uuid.UUID  `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	LevelID         *uuid.UUID `gorm:"column:elegibility_level_id;type:varchar(36);"`
	Level           *Level     `gorm:"foreignKey:LevelID;constraint:OnDelete:SET NULL;"`
	CustomerID      *uuid.UUID `gorm:"column:customer_id;type:varchar(36);uniqueIndex:idx_unique;"`
	Customer        *Customer  `gorm:"foreignKey:CustomerID;constraint:OnDelete:SET NULL;"`
	ValidityMonth   *int       `gorm:"column:validity_month;type:TINYINT;not null;default:1;check:validity_month > -1;uniqueIndex:idx_unique;"`
	ValidityYear    *int       `gorm:"column:validity_year;type:SMALLINT;not null;default:1;check:validity_year > -1;uniqueIndex:idx_unique;"`
	ManuallyTouched *bool      `gorm:"column:manually_touched;type:boolean;not null;default:false;"`
	CreatedByID     *uuid.UUID `gorm:"column:created_by_id;type:varchar(36);"`
	CreatedBy       *User      `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
	UpdatedByID     *uuid.UUID `gorm:"column:updated_by_id;type:varchar(36);"`
	UpdatedBy       *User      `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:SET NULL;"`

	gorm.Model
}

func (ul *CustomerLevel) TableName() string {
	return "customer_levels"
}

func (ul *CustomerLevel) BeforeCreate(tx *gorm.DB) (err error) {
	ul.ID = uuid.New()

	return
}

func (l *Level) TableName() string {
	return "levels"
}

func (l *Level) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()

	return
}
