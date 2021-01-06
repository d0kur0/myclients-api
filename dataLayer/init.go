package dataLayer

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() (err error) {
	db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return
	}

	if err = db.AutoMigrate(&User{}, &AuthToken{}); err != nil {
		return
	}

	return
}

func GetDB() *gorm.DB {
	return db
}
