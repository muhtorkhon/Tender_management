package models

import "time"

type Offers struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TenderID     uint       `gorm:"not null" json:"tender_id" binding:"required"`
	ContractorID uint       `gorm:"not null" json:"contractor_id" binding:"required"`
	Price        float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	DeliveryTime *time.Time `gorm:"not null" json:"delivery_time"`
	Comments     string     `gorm:"type:text;not null" json:"comments"`
	Status       bool       `gorm:"default:true" json:"status"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
	CreatedAt    *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Users        *Users     `gorm:"foreignKey:ContractorID" json:"-"`
	Tenders      *Tenders   `gorm:"foreignKey:TenderID" json:"-"`
}

type OffersRequest struct {
	TenderID     uint    `json:"tender_id" binding:"required"`
	ContractorID uint    `json:"contractor_id" binding:"required"`
	Price        float64 `json:"price" binding:"required,gt=0"`
	DeliveryTime string  `json:"delivery_time" binding:"required"`
	Comments     string  `json:"comments" binding:"required"`
	Status       bool    `json:"status"`
}

type Stats struct {
	MinPrice     float64   `json:"min_price"`
	MaxPrice     float64   `json:"max_price"`
	MinDelivery  time.Time `json:"min_delivery"`
	MaxDelivery  time.Time `json:"max_delivery"`
	TotalRecords int64     `json:"total_records"`
}

