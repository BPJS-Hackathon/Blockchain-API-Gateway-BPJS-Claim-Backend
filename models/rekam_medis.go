package models

const (
	RawatJalan = "RAWAT JALAN"
	RawatInap  = "RAWAT INAP"

	OutcomeSembuh      = "SEMBUH"
	OutcomeRujuk       = "RUJUK"
	OutcomePulangPaksa = "PULANG PAKSA"
	OutcomeMeninggal   = "MENINGGAL"
)

type DiagnosisCode struct {
	Code       string `gorm:"type:varchar(10);primaryKey" json:"code"`
	Keterangan string `gorm:"type:varchar(255);not null" json:"keterangan"`
	Harga      uint   `gorm:"type:int;not null" json:"harga"`
}

type RekamMedis struct {
	ID            string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PesertaNIK    string        `gorm:"type:varchar(20);not null" json:"peserta_nik"`
	UserID        string        `gorm:"type:uuid;not null" json:"user_id"`
	DiagnosisCode string        `gorm:"type:varchar(10);not null" json:"diagnosis_code"`
	Note          string        `gorm:"type:text" json:"note"`
	JenisRawat    string        `gorm:"type:varchar(50);not null" json:"jenis_rawat"`
	AdmissionDate int64         `gorm:"type:bigint;not null" json:"admission_date"`
	DischargeDate *int64        `gorm:"type:bigint" json:"discharge_date"`
	Outcome       *string       `gorm:"type:varchar(100)" json:"outcome"`
	User          User          `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Diagnosis     DiagnosisCode `gorm:"foreignKey:DiagnosisCode;references:Code" json:"diagnosis"`
	Claims        []Claims      `gorm:"foreignKey:RekamMedisID;references:ID" json:"claims"`
}
