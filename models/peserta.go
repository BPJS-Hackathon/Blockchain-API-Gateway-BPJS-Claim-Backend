package models

import (
	"time"
)

type PesertaBPJS struct {
	// PSTV01 - Nomor peserta (Primary Key)
	PSTV01 string `gorm:"column:pstv01;primaryKey;type:varchar(32)" json:"nomor_peserta"`
	// PSTV02 - Nomor keluarga
	PSTV02 string `gorm:"column:pstv02;index;type:varchar(32)" json:"nomor_keluarga"`
	// PSTV03 - Tanggal lahir
	PSTV03 time.Time `gorm:"column:pstv03;type:date" json:"tanggal_lahir"`
	// PSTV04 - Hubungan keluarga (1: Peserta, 2: Suami, 3: Istri, 4: Anak...)
	PSTV04 int `gorm:"column:pstv04;type:smallint" json:"hubungan_keluarga"`
	// PSTV05 - Jenis kelamin (1: Laki-laki, 2: Perempuan)
	PSTV05 int `gorm:"column:pstv05;type:smallint" json:"jenis_kelamin"`
	// PSTV06 - Status perkawinan (1: Belum kawin, 2: Kawin, 3: Cerai)
	PSTV06 int `gorm:"column:pstv06;type:smallint" json:"status_perkawinan"`
	// PSTV07 - Kelas rawat (I / II / III)
	PSTV07 string `gorm:"column:pstv07;type:varchar(10)" json:"kelas_rawat"`
	// PSTV08 - Segmentasi peserta (1: BP, 2: PBI APBN, 3: PBI APBD...)
	PSTV08 int `gorm:"column:pstv08;type:smallint" json:"segmentasi_peserta"`
	// PSTV09 - Provinsi tempat tinggal (kode BPS)
	PSTV09 string `gorm:"column:pstv09;type:varchar(4)" json:"provinsi_tinggal"`
	// PSTV10 - Kabupaten/Kota tempat tinggal
	PSTV10 string `gorm:"column:pstv10;type:varchar(4)" json:"kabupaten_kota_tinggal"`
	// PSTV11 - Kepemilikan faskes
	PSTV11 string `gorm:"column:pstv11;type:varchar(64)" json:"kepemilikan_faskes"`
	// PSTV12 - Jenis fasilitas kesehatan (1: Puskesmas, 2: Klinik Pratama...)
	PSTV12 int `gorm:"column:pstv12;type:smallint" json:"jenis_faskes"`
	// PSTV13 - Provinsi faskes
	PSTV13 string `gorm:"column:pstv13;type:varchar(4)" json:"provinsi_faskes"`
	// PSTV14 - Kabupaten/Kota faskes
	PSTV14 string `gorm:"column:pstv14;type:varchar(4)" json:"kabupaten_kota_faskes"`
	// PSTV15 - Bobot (faktor pengali)
	PSTV15 float64 `gorm:"column:pstv15;type:decimal(10,4)" json:"bobot"`
	// PSTV16 - Tahun sampel
	PSTV16 int `gorm:"column:pstv16;type:integer" json:"tahun_sampel"`
	// PSTV17 - Status kepesertaan
	PSTV17 string `gorm:"column:pstv17;type:varchar(16)" json:"status_kepesertaan"`
	// PSTV18 - Tahun meninggal (nullable)
	PSTV18    *int   `gorm:"column:pstv18;type:integer" json:"tahun_meninggal"`
	CreatedBy string `gorm:"column:created_by;type:uuid" json:"created_by"`
}

func (PesertaBPJS) TableName() string {
	return "peserta_bpjs"
}
