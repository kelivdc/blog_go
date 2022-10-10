package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	CategoryID  int      `json:"category_id"`
	Category    Category `json:"category"`
	Title       string   `json:"title" validate:"required" gorm:"size:100"`
	Body        string   `json:"body" validate:"required"`
	ShortDesc   string   `json:"short_desc" validate:"required" gorm:"size:50"`
	Description string   `json:"description" validate:"required" gorm:"size:100"`
	Keyword     string   `json:"keyword" gorm:"size:100"`
	Slug        string   `json:"slug" validate:"required,min=3,max=200"`
	Image       string   `json:"image" validate:"required" gorm:"size:100"`
	ImageNote   string   `json:"image_note" gorm:"size:50"`
	Publish     bool     `json:"publish" validate:"required"`
	gorm.Model  `json:"-"`
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("Created")
	return
}
