package models

import "github.com/golang-jwt/jwt/v5"

type Users struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstName   string `gorm:"type:varchar(255);not null" json:"first_name"`
	Email       string `gorm:"type:varchar(255);unique;not null" json:"email"`
	PhoneNumber string `gorm:"type:varchar(255);unique;not null" json:"phone_number"`
	Password    string `gorm:"type:varchar(255);not null" json:"password"`
	Role        string `gorm:"type:varchar(255);not null" json:"user_role"`
	IsActive    bool   `gorm:"default:false" json:"is_active"`
}

type UserRegister struct {
	FirstName   string `json:"first_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Password    string `json:"password" validate:"password"`
	Role        string `json:"role" binding:"required,oneof=client contractor"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type ResetPassword struct {
	ConfirmPassword string `json:"confirm_password"`
	NewPassword     string `json:"new_password"`
	UserID          uint   `json:"user_id"`
}

type ForgotPassword struct {
	PhoneNumber string `json:"phone_number"`
}

type NewPassword struct {
	UserID      uint   `json:"user_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

