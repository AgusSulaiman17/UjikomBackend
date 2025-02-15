package models

import (
	"time"
)

// Penulis represents the penulis table in the database
type Penulis struct {
	IDPenulis  uint      `json:"id_penulis" gorm:"primaryKey"`
	Nama       string    `json:"nama" gorm:"not null"`
	DibuatPada time.Time `json:"dibuat_pada" gorm:"autoCreateTime"`
	DiperbaruiPada time.Time `json:"diperbarui_pada" gorm:"autoUpdateTime"`
}
