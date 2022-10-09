package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model    `json:"-"`
	ID            uint      `json:"id" gorm:"primaryKey"`
	CategoryRefer int       `json:"category_id" validate:"required"`
	Category      Category  `json:"-" gorm:"foreignKey:CategoryRefer"`
	Title         string    `json:"title" validate:"required" gorm:"size:100"`
	Body          string    `json:"body" validate:"required"`
	ShortDesc     string    `json:"short_desc" validate:"required" gorm:"size:50"`
	Description   string    `json:"description" validate:"required" gorm:"size:100"`
	Keyword       string    `json:"keyword" gorm:"size:100"`
	Slug          string    `json:"slug" validate:"required,min=3,max=200"`
	Image         string    `json:"image" validate:"required" gorm:"size:100"`
	ImageNote     string    `json:"image_note" gorm:"size:50"`
	Publish       bool      `json:"publish" validate:"required"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("Created")
	return
}
