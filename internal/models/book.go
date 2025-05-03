package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `gorm:"size:255;not null" json:"title"`
	Author    string `gorm:"size:255;not null" json:"author"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Category  string `gorm:"size:255" json:"category"`
	ImageUrl  string `gorm:"size:255" json:"image_url"`
	IsPrivate bool   `json:"is_private" gorm:"default:false"`
}
