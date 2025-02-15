package models

import (
)

type User struct {
	IDUser          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name" binding:"required"`
	Email       string `gorm:"uniqueIndex" json:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" binding:"omitempty"`
	Role        string `gorm:"default:user" json:"role"` 
	Image       string `json:"image"`                   
	Alamat      string `json:"alamat"`                  
	NoTelepon   string `json:"no_telepon"`             
}