package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campaign struct {
	ID          uuid.UUID      `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	Name        string         `gorm:"column:name;type:varchar(255);not null;default:'';"`
	Discount    *float64       `gorm:"column:discount;type:float;not null;default:0;check:discount > -1;"`
	ValidFrom   *time.Time     `gorm:"column:valid_from;type:datetime;not null;default:NOW();"`
	ValidTo     *time.Time     `gorm:"column:valid_to;type:datetime;not null;default:NOW();"`
	Active      *bool          `gorm:"column:active;type:boolean;not null;default:true;"`
	GasStations *[]*GasStation `gorm:"many2many:gas_stations_campaigns;"`
	CreatedByID *uuid.UUID     `gorm:"column:created_by_id;type:varchar(36);"`
	CreatedBy   *User          `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
	UpdatedByID *uuid.UUID     `gorm:"column:updated_by_id;type:varchar(36);"`
	UpdatedBy   *User          `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:SET NULL;"`

	gorm.Model
}

func (c *Campaign) TableName() string {
	return "campaigns"
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()

	return
}

// Copier helper magic method in order to get the right date
func (c *Campaign) ValidToStr(validTo string) {
	toDate, _ := time.Parse("2006-01-02 15:04:05", validTo)

	c.ValidTo = &toDate
}

// Copier helper magic method in order to get the right date
func (c *Campaign) ValidFromStr(validFrom string) {
	fromDate, _ := time.Parse("2006-01-02 15:04:05", validFrom)

	c.ValidFrom = &fromDate
}

// Copier helper magic method in order to get the right date
func (c *Campaign) ValidToStrUpdate(validTo *string) {
	toDate, _ := time.Parse("2006-01-02 15:04:05", *validTo)

	c.ValidTo = &toDate
}

// Copier helper magic method in order to get the right date
func (c *Campaign) ValidFromStrUpdate(validFrom *string) {
	fromDate, _ := time.Parse("2006-01-02 15:04:05", *validFrom)

	c.ValidFrom = &fromDate
}
