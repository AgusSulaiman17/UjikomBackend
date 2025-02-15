package models

import "time"

type Buku struct {
	IDBuku        uint       `json:"id_buku" gorm:"primaryKey"`
	Judul         string     `json:"judul" gorm:"not null"`
	IDPenerbit    uint       `json:"id_penerbit"`
	Penerbit      Penerbit   `json:"penerbit" gorm:"foreignKey:IDPenerbit;references:IDPenerbit"`
	IDPenulis     uint       `json:"id_penulis"`
	Penulis       Penulis    `json:"penulis" gorm:"foreignKey:IDPenulis;references:IDPenulis"`
	IDKategori    uint       `json:"id_kategori"`
	Kategori      Kategori   `json:"kategori" gorm:"foreignKey:IDKategori;references:IDKategori"`
	Deskripsi     string     `json:"deskripsi"`
	Jumlah        int        `json:"jumlah"`
	Gambar        string     `json:"gambar"`
	Status        bool       `json:"status" gorm:"default:true"`
	ISBN          string     `json:"isbn" gorm:"type:varchar(20);not null;unique"`
	DibuatPada    time.Time  `json:"dibuat_pada" gorm:"autoCreateTime"`
	DiperbaruiPada time.Time `json:"diperbarui_pada" gorm:"autoUpdateTime"`
}

