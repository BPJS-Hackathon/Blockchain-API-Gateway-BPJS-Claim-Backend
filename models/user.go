package models

const (
	RoleFaskes  = "faskes"
	RoleAdmin   = "admin"
	RoleAuditor = "auditor"

	JenisFaskes1 = "FASKES 1"
	JenisFaskes2 = "FASKES 2"
)

type User struct {
	ID           string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string  `gorm:"uniqueIndex;type:varchar(64);not null" json:"username"`
	PasswordHash string  `gorm:"type:varchar(255);not null" json:"-"`
	Role         string  `gorm:"type:varchar(32);not null" json:"role"` // admin, faskes, auditor
	FaskesID     *string `gorm:"type:uuid" json:"faskes_id,omitempty"`
	Faskes       *Faskes `gorm:"foreignKey:FaskesID;references:ID" json:"faskes,omitempty"`
}

// models/faskes.go

type Faskes struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NamaFaskes  string `gorm:"type:varchar(100);not null" json:"nama_faskes"`
	JenisFaskes string `gorm:"type:varchar(100);not null" json:"jenis_faskes"`
}

type Faskes1Repo interface {
	CreateRekamMedis(rm RekamMedis) error
}

type Faskes1Service interface {
	CreateRekamMedis(rm RekamMedis) error
}

type Faskes2Repo interface {
	GetAllDiagnosisCodes() ([]DiagnosisCode, error)
	CreateRekamMedisandClaim(rm RekamMedis, cl Claims) (string, error)
}

type Faskes2Service interface {
	GetAllDiagnosisCodes() ([]DiagnosisCode, error)
	CreateRekamMedisandClaim(rm RekamMedis, cl Claims) (string, error)
}
