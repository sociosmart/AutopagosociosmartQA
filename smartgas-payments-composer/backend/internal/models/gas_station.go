package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GasStation struct {
	ID            uuid.UUID   `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	ExternalID    string      `gorm:"column:external_id;type:varchar(10);not null;unique;"`
	Name          string      `gorm:"column:name;type:varchar(255);not null;default:'';uniqueIndex:idx_name_ip;"`
	Ip            string      `gorm:"column:ip;type:varchar(15);not null;default:'';uniqueIndex:idx_name_ip;"`
	CrePermission string      `gorm:"column:cre_permission;type:varchar(150);not null;default:'';unique;"`
	Street        string      `gorm:"column:street;type:varchar(150);not null;default:'';"`
	ZipCode       string      `gorm:"column:zip_code;type:varchar(8);not null;default:'';"`
	City          string      `gorm:"column:city;type:varchar(50);not null;default:'';"`
	State         string      `gorm:"column:state;type:varchar(50);not null;default:'';"`
	Neighborhood  string      `gorm:"column:neighborhood;type:varchar(50);not null;default:'';"`
	OutsideNumber string      `gorm:"column:outside_number;type:varchar(10);not null;default:'';"`
	Latitude      string      `gorm:"column:latitude;type:varchar(50);not null;default:'';"`
	Longitude     string      `gorm:"column:longitude;type:varchar(50);not null;default:'';"`
	Active        *bool       `gorm:"column:active;type:boolean;not null;default:true;"`
	LegalNameID   string      `gorm:"column:legal_name_id;type:varchar(10);not null;default:'';"`
	CreatedByID   *uuid.UUID  `gorm:"column:created_by_id;type:varchar(36);"`
	CreatedBy     *User       `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;"`
	UpdatedByID   *uuid.UUID  `gorm:"column:updated_by_id;type:varchar(36);"`
	UpdatedBy     *User       `gorm:"foreignKey:UpdatedByID;constraint:OnDelete:SET NULL;"`
	Users         []*User     `gorm:"many2many:user_gas_stations;"`
	Campaigns     []*Campaign `gorm:"many2many:gas_stations_campaigns;"`
	gorm.Model
}

func (gs *GasStation) TableName() string {
	return "gas_stations"
}

func (gs *GasStation) BeforeCreate(tx *gorm.DB) (err error) {
	gs.ID = uuid.New()

	return
}
