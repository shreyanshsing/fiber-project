package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex; not null; <-:create" json:"user_name" validate:"required"`
	Email    string `gorm:"uniqueIndex; not null; <-:create" json:"email" validate:"required,email"`
	Password string `gorm:"not null; <-" json:"password" validate:"required"`
	Names    string `json:"names"`
	Verified bool   `gorm:"default:false" json:"verified"`
}
