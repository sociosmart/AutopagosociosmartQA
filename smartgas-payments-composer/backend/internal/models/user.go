package models

import (
	"errors"
	"net/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID     `gorm:"column:id;primaryKey;type:varchar(36);<-:create;"`
	FirstName   string        `gorm:"column:first_name;type:varchar(100);not null;"`
	LastName    string        `gorm:"column:last_name;type:varchar(100);default:''; not null;"`
	Email       string        `gorm:"column:email;type:varchar(255);uniqueIndex;not null;"`
	Password    string        `gorm:"column:password;type:varchar(255);not null;"`
	Active      *bool         `gorm:"column:active;type:boolean;not null;default:true;"`
	IsAdmin     *bool         `gorm:"column:is_admin;type:boolean;not null;default:false;"`
	Groups      []*Group      `gorm:"many2many:user_groups;"`
	Permissions []*Permission `gorm:"many2many:user_permissions;"`
	GasStations []*GasStation `gorm:"many2many:user_gas_stations;"`
	gorm.Model
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	if _, err = mail.ParseAddress(u.Email); err != nil {
		return errors.New("'email' field is not a valid email")
	}
	u.ID = uuid.New()
	u.HashPassword()

	return
}

func (u *User) HashPassword() (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	u.Password = string(hashedPassword)

	return
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {

	if u.Password != "" {
		u.HashPassword()
	}

	return
}
