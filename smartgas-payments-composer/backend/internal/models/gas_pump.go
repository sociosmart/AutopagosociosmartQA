package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GasPump struct {
	ID           uuid.UUID   `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ExternalID   string      `gorm:"column:external_id;type:varchar(10);not null;unique;"`
	RegularPrice *float64    `gorm:"column:regular_price;type:double;not null;default:0;check:regular_price > -1;"`
	PremiumPrice *float64    `gorm:"column:premium_price;type:double;not null;default:0;check:premium_price > -1;"`
	DieselPrice  *float64    `gorm:"column:diesel_price;type:double;not null;default:0;check:diesel_price > -1;"`
	Number       string      `gorm:"column:number;type:varchar(3);default:'01';not null;uniqueIndex:idx_gas_station_number;"`
	GasStationID *uuid.UUID  `gorm:"column:gas_station_id;type:varchar(36);uniqueIndex:idx_gas_station_number;"`
	GasStation   *GasStation `gorm:"constraint:OnDelete:SET NULL;"`
	Active       *bool       `gorm:"column:active;type:boolean;not null;default:true;"`
	CreatedByID  *uuid.UUID  `gorm:"column:created_by_id;type:varchar(36);"`
	CreatedBy    *User       `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
	UpdatedByID  *uuid.UUID  `gorm:"column:updated_by_id;type:varchar(36);"`
	UpdatedBy    *User       `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:SET NULL;"`
	gorm.Model
}

func (gp *GasPump) TableName() string {
	return "gas_pumps"
}

func (gp *GasPump) BeforeCreate(tx *gorm.DB) (err error) {
	gp.ID = uuid.New()

	return
}
