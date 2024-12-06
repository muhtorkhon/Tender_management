package models

import "time"

type Notif struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	Message    string     `gorm:"type:text;not null" json:"message"`
	RelationID uint       `gorm:"not null" json:"relation_id"`
	Type       string     `gorm:"type:varchar(50);not null" json:"type"`
	CreatedAt  *time.Time `gorm:"autoCreateTime" json:"created_at"`
	Users      *Users     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	Tenders    *Tenders   `gorm:"foreignKey:RelationID;constraint:OnDelete:CASCADE;" json:"-"`
	Offers     *Offers    `gorm:"foreignKey:RelationID;constraint:OnDelete:CASCADE;" json:"-"`
}

type NotifRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
	RelationID uint   `json:"relation_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
}
