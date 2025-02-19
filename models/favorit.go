package models

import "time"

type Favorit struct {
    IDFavorit uint      `json:"id" gorm:"primaryKey"`
    IDUser    uint      `json:"id_user"`
    IDBuku    uint      `json:"id_buku"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    
    User User `gorm:"foreignKey:IDUser;references:IDUser" json:"user"`
    Buku Buku `gorm:"foreignKey:IDBuku;references:IDBuku" json:"buku"`
}
