package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" validate:"required" gorm:"size:50"`
	Slug       string    `json:"slug" validate:"required,min=3,max=100"`
	Publish    bool      `json:"publish"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
