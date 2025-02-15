package models

import "time"

type Kategori struct {
	IDKategori   uint      `gorm:"primaryKey" json:"id"`
	Kategori     string    `gorm:"type:varchar(100);not null" json:"kategori"`
	DibuatPada   time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}
