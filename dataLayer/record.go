package dataLayer

import "time"

type Record struct {
	Model
	UserID    uint64    `json:"userId"`
	User      User      `json:"user"`
	ServiceID uint64    `json:"serviceId"`
	Service   Service   `json:"service"`
	ClientID  uint64    `json:"clientId"`
	Client    Client    `json:"client"`
	Date      time.Time `json:"date"`
}
