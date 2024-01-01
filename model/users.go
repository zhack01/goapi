package model

import (
	"time"

	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"DEFAULT:NULL"`
	Username       string `gorm:"size:45,unique"`
	Email          string `gorm:"size:255"`
	Password       string
	PasswordString string
	UserType       string
	IsAdmin        string `gorm:"DEFAULT:2"`
	OperatorId     int
	ClientId       int `gorm:"DEFAULT:0"`
	BrandId        int `gorm:"DEFAULT:0"`
	StatusId       int `json:"status_id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
