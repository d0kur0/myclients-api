package dataLayer

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() (err error) {
	db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return
	}

	if err = db.AutoMigrate(&User{}, &AuthToken{}, &Client{}, &Service{}, &Record{}); err != nil {
		return
	}

	return
}

func GetDB() *gorm.DB {
	return db
}
