package dataLayer

import (
	"time"
)

type User struct {
	Model
	Name           string      `json:"name"`
	Email          string      `json:"email" gorm:"uniqueIndex"`
	Password       string      `json:"password"`
	AvatarPath     string      `json:"avatarPath"`
	IsEmailConfirm bool        `json:"isEmailConfirm"`
	AuthTokens     []AuthToken `json:"authTokens"`
}

type AuthToken struct {
	Model
	UserID   uint64    `json:"userId"`
	Token    string    `json:"token"`
	DeadTime time.Time `json:"deadTime"`
}
