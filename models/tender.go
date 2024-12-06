package models

import (
	"time"
)

type Tenders struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Deadline    *time.Time `gorm:"not null" json:"deadline"`
	Budget      float64    `gorm:"type:decimal(10,2);not null" json:"budget"`
	FileURL     string     `gorm:"type:varchar(255)" json:"file_url,omitempty"`
	Status      bool       `gorm:"default:true" json:"status"`
	ClientID    uint       `gorm:"not null" json:"client_id" binding:"required"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Users       *Users     `gorm:"foreignKey:ClientID" json:"-"`
}

type TenderRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Deadline    string  `json:"deadline" binding:"required"`
	Budget      float64 `json:"budget" binding:"required,gt=0"`
	FileURL     string  `json:"file_url,omitempty"`
	ClientID    uint    `json:"client_id" binding:"required"`
}
