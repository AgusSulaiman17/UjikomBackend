package models

import (
	"time"
)

type Penerbit struct {
	IDPenerbit    int       `json:"id_penerbit" gorm:"primaryKey;autoIncrement"`
	Nama          string    `json:"nama" gorm:"not null" validate:"required"`
	DibuatPada    time.Time `json:"dibuat_pada" gorm:"default:CURRENT_TIMESTAMP"`
	DiperbaruiPada time.Time `json:"diperbarui_pada" gorm:"default:CURRENT_TIMESTAMP"`
}
