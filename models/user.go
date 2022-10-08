package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	Email      string `json:"email" validate:"required,email,min=6,max=100" gorm:"unique"`
	Password   string `json:"password" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Active     bool   `json:"active" validate:"required"`
	Token      string `json:"token"`
}
