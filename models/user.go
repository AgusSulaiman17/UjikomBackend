package models

type User struct {
	IDUser    uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name" form:"name" binding:"required"`
	Email     string `gorm:"uniqueIndex" json:"email" form:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" form:"password" binding:"omitempty"`
	Role      string `gorm:"default:user" json:"role" form:"role"`
	Image     string `json:"image" form:"image"`
	Alamat    string `json:"alamat" form:"alamat"`
	NoTelepon string `json:"no_telepon" form:"no_telepon"`
}
