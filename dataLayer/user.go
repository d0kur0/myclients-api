package dataLayer

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             uint64      `json:"id"`
	Name           string      `json:"name"`
	Email          string      `json:"email" gorm:"uniqueIndex"`
	Password       string      `json:"password"`
	AvatarPath     string      `json:"avatarPath"`
	IsEmailConfirm bool        `json:"isEmailConfirm"`
	AuthTokens     []AuthToken `json:"authTokens"`
}

// Yes, i'm gay
type AuthToken struct {
	gorm.Model
	UserID uint64 `json:"userId"`
	Token  string `json:"token"`
}
