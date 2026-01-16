package database

import (
	"fmt"
	"smartgas-payment/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	fmt.Println("Applying migrations...")
	if err := db.AutoMigrate(
		models.User{},
		models.GasStation{},
		models.GasPump{},
		models.Customer{},
		models.Payment{},
		models.Synchronization{},
		models.SynchronizationDetail{},
		models.SynchronizationError{},
		models.PaymentEvent{},
		models.AuthorizedApplication{},
		models.Permission{},
		models.Group{},
		models.Setting{},
		models.Campaign{},
		models.Level{},
		models.CustomerLevel{},
	); err != nil {
		panic(err)
	}
	fmt.Println("Migrations applied...")
}
