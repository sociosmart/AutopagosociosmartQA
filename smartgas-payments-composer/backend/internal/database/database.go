package database

import (
	"fmt"
	"log"
	"smartgas-payment/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(c config.Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", c.DB.User, c.DB.Pass, c.DB.Host, c.DB.Port, c.DB.Name)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			break
		}

		time.Sleep(time.Second * 1)
		log.Println("retrying to connect to database...")
	}

	return

}

func CloseConnection(db *gorm.DB) {

	if db != nil {
		dbB, _ := db.DB()

		dbB.Close()

	}

}
