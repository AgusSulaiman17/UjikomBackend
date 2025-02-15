package models

import "time"

type Peminjaman struct {
    IDPeminjaman   uint      `gorm:"primaryKey" json:"id_peminjaman"`
    IDUser         uint      `json:"id_user"`
    User           User      `gorm:"foreignKey:IDUser;references:IDUser" json:"user"`
    IDBuku         uint      `json:"id_buku"`
    Buku           Buku      `gorm:"foreignKey:IDBuku;references:IDBuku" json:"buku"`
    TanggalPinjam  time.Time `json:"tanggal_pinjam"`
    DurasiHari     int       `json:"durasi_hari"`
    TanggalKembali time.Time `json:"tanggal_kembali"`
    StatusKembali  bool      `json:"status_kembali"`
    Denda          float64   `json:"denda"`
    DibuatPada     time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
    DiperbaruiPada time.Time `gorm:"autoUpdateTime" json:"diperbarui_pada"`
}
