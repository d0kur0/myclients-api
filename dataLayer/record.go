package dataLayer

import "time"

type Record struct {
	Model
	UserID   uint64    `json:"userId"`
	User     User      `json:"user"`
	Services []Service `json:"services" gorm:"many2many:record_services"`
	ClientID uint64    `json:"clientId"`
	Client   Client    `json:"client"`
	Date     time.Time `json:"date"`
}
