package models

import "time"

type Peminjaman struct {
    IDPeminjaman   uint      `gorm:"primaryKey" json:"id_peminjaman"`
    IDUser         uint      `gorm:"column:id_user" json:"id_user"`
    User           User      `gorm:"foreignKey:IDUser;references:IDUser" json:"user,omitempty"`
    IDBuku         uint      `gorm:"column:id_buku" json:"id_buku"`
    Buku           Buku      `gorm:"foreignKey:IDBuku;references:IDBuku" json:"buku,omitempty"`
    TanggalPinjam  time.Time `json:"tanggal_pinjam"`
    DurasiHari     int       `json:"durasi_hari"`
    TanggalKembali time.Time `json:"tanggal_kembali"`
    Status         string    `json:"status"`
    StatusKembali  bool      `json:"status_kembali"`
    Denda          float64   `json:"denda"`
    DibuatPada     time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
    DiperbaruiPada time.Time `gorm:"autoUpdateTime" json:"diperbarui_pada"`
    IsDeletedByUser bool `gorm:"default:false" json:"deleted_by_user"`
}

